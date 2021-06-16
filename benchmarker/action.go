package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"
)

func BrowserAccess(ctx context.Context, step *isucandar.BenchmarkStep, user *User, rpath string) error {
	req, err := user.Agent.GET(rpath)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return err
	}

	if err := assertStatusCode(res, 200); err != nil {
		if res.StatusCode != 304 {
			return err
		}
	}

	resources, perr := user.Agent.ProcessHTML(ctx, res, res.Body)
	if perr != nil {
		return failure.NewError(ErrCritical, err)
	}

	resStatusCodes := map[string]int{
		"/esm/index.js": 200,
	}

	for _, resource := range resources {
		if resource.Error != nil {
			step.AddError(failure.NewError(ErrInvalidAsset, resource.Error))
			continue
		}

		if resource.Response.StatusCode == 304 {
			continue
		}

		if statusCode, ok := resStatusCodes[resource.Request.URL.Path]; ok {
			if err := assertStatusCode(resource.Response, statusCode); err != nil {
				step.AddError(failure.NewError(ErrInvalidAsset, err))
			}
		}

		if err := assertChecksum(resource.Response); err != nil {
			step.AddError(err)
		}
	}

	return nil
}

type SignupResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

func ActionSignup(ctx context.Context, step *isucandar.BenchmarkStep, u *User) error {
	values := url.Values{}
	values.Add("email", u.Email)
	values.Add("nickname", u.Nickname)

	body := strings.NewReader(values.Encode())

	req, err := u.Agent.POST("/api/signup", body)
	if err != nil {
		// request が生成できないなんてのは相当やばい状況なのでたいてい Critical です
		// さっさと Critical エラーにして早めにベンチマーカー止めてあげるのも優しさ
		return failure.NewError(ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := u.Agent.Do(ctx, req)
	if err != nil {
		return err
	}

	hasError := false
	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
		hasError = true
	}

	if err := assertContentType(res, "application/json"); err != nil {
		step.AddError(err)
		hasError = true
	}

	jsonResp := &SignupResponse{}
	if err := assertJSONBody(res, jsonResp); err != nil {
		step.AddError(err)
		hasError = true
	} else {
		if err := assertEqualString(u.Email, jsonResp.Email, "signup.email"); err != nil {
			step.AddError(err)
			hasError = true
		}

		if err := assertEqualString(u.Nickname, jsonResp.Nickname, "signup.nickname"); err != nil {
			step.AddError(err)
			hasError = true
		}
	}

	if !hasError {
		u.ID = jsonResp.ID
		u.CreatedAt = jsonResp.CreatedAt
		step.AddScore(ScoreSignup)
	}

	return nil
}

// ユーザーをたくさんつくるよ
func ActionSignups(parent context.Context, step *isucandar.BenchmarkStep, s *Scenario) error {
	// とりあえず10秒くらい
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// とりあえず50並列くらい
	w, err := worker.NewWorker(func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			// context が終わってたら抜ける
			// あ、Paralle だと一回しか実行しないのか
			return
		default:
		}

		wg.Add(1)
		defer wg.Done()

		user, err := s.NewUser()
		if err != nil {
			step.AddError(err)
			return
		}
		if err := ActionSignup(parent, step, user); err != nil {
			step.AddError(err)
			return
		}
		s.Users.Add(user)
	}, worker.WithMaxParallelism(s.Parallelism), worker.WithInfinityLoop())
	if err != nil {
		return err
	}

	// 一応ここでも待ち合わせはするんだけどね
	w.Process(ctx)

	// 確実に止める、止まったことを検知するために
	wg.Done()
	wg.Wait()

	return nil
}

type LoginResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

// Action がエラーを返す → Action の失敗
// Action がエラーを返さない → Action としては成功。シナリオとしてはどうかわからない
func ActionLogin(ctx context.Context, step *isucandar.BenchmarkStep, u *User) error {
	values := url.Values{}

	if u.FailOnLogin {
		values.Add("email", "invalid-"+u.Email)
	} else {
		values.Add("email", u.Email)
	}

	body := strings.NewReader(values.Encode())

	req, err := u.Agent.POST("/api/login", body)
	if err != nil {
		// request が生成できないなんてのは相当やばい状況なのでたいてい Critical です
		// さっさと Critical エラーにして早めにベンチマーカー止めてあげるのも優しさ
		return failure.NewError(ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := u.Agent.Do(ctx, req)
	if err != nil {
		return err
	}

	hasError := false
	if u.FailOnLogin {
		if err := assertStatusCode(res, 403); err != nil {
			step.AddError(err)
			hasError = true
		}
	} else {
		if err := assertStatusCode(res, 200); err != nil {
			step.AddError(err)
			hasError = true
		}

		if err := assertContentType(res, "application/json"); err != nil {
			step.AddError(err)
			hasError = true
		}

		jsonResp := &LoginResponse{}
		if err := assertJSONBody(res, jsonResp); err != nil {
			step.AddError(err)
			hasError = true
		} else {
			if err := assertEqualString(u.Email, jsonResp.Email, "login.email"); err != nil {
				step.AddError(err)
				hasError = true
			}
			if err := assertEqualString(u.Nickname, jsonResp.Nickname, "login.nickname"); err != nil {
				step.AddError(err)
				hasError = true
			}
		}
	}

	if !hasError {
		step.AddScore(ScoreLogin)
	}

	return nil
}

func ActionLogins(parent context.Context, step *isucandar.BenchmarkStep, s *Scenario) error {
	usersCount := s.Users.Count()

	// とりあえず30秒耐える
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// とりあえず100並列くらい
	w, err := worker.NewWorker(func(ctx context.Context, idx int) {
		select {
		case <-ctx.Done():
			// context が終わってたら抜ける
			return
		default:
		}

		user := s.Users.Get(idx)
		if err := ActionLogin(ctx, step, user); err != nil {
			step.AddError(err)
		}
	}, worker.WithMaxParallelism(s.Parallelism), worker.WithLoopCount(int32(usersCount)))
	if err != nil {
		return err
	}

	w.Process(ctx)

	wg.Done()
	wg.Wait()

	return nil
}

type CreateScheduleResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Capacity  uint      `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}

func ActionCreateSchedule(ctx context.Context, step *isucandar.BenchmarkStep, s *Scenario) (*Schedule, error) {
	user := s.StaffUser

	schedule, err := s.NewSchedule()
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("title", schedule.Title)
	values.Add("capacity", strconv.Itoa(int(schedule.Capacity)))

	body := strings.NewReader(values.Encode())
	req, err := user.Agent.POST("/api/schedules", body)
	if err != nil {
		return nil, failure.NewError(ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if !user.Staff {
		if err := assertStatusCode(res, 401); err != nil {
			return nil, err
		}
	}

	hasError := false

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
		hasError = true
	}

	if err := assertContentType(res, "application/json"); err != nil {
		step.AddError(err)
		hasError = true
	}

	jsonResp := &CreateScheduleResponse{}
	if err := assertJSONBody(res, jsonResp); err != nil {
		step.AddError(err)
		hasError = true
	} else {
		if err := assertEqualString(jsonResp.Title, schedule.Title, "create-schedule.title"); err != nil {
			step.AddError(err)
			hasError = true
		}
		if err := assertEqualUint(jsonResp.Capacity, schedule.Capacity, "create-schedule.capacity"); err != nil {
			step.AddError(err)
			hasError = true
		}
	}

	if !hasError {
		schedule.ID = jsonResp.ID
		schedule.CreatedAt = jsonResp.CreatedAt

		step.AddScore(ScoreCreateSchedule)
	}
	return schedule, nil
}

type SchedulesResponseItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Capacity uint   `json:"capacity"`
	Reserved uint   `json:"reserved"`
}

func ActionGetSchedules(ctx context.Context, step *isucandar.BenchmarkStep, user *User) ([]*SchedulesResponseItem, error) {
	req, err := user.Agent.GET("/api/schedules")
	if err != nil {
		return nil, failure.NewError(ErrCritical, err)
	}

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
	}

	if err := assertContentType(res, "application/json"); err != nil {
		step.AddError(err)
	}

	schedules := []*SchedulesResponseItem{}
	if err := assertJSONBody(res, &schedules); err != nil {
		step.AddError(err)
	}

	return schedules, nil
}

func ActionGetAllSchedules(ctx context.Context, step *isucandar.BenchmarkStep, user *User) ([]*SchedulesResponseItem, error) {
	req, err := user.Agent.GET("/api/schedules?reserved=1")
	if err != nil {
		return nil, failure.NewError(ErrCritical, err)
	}

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
	}

	if err := assertContentType(res, "application/json"); err != nil {
		step.AddError(err)
	}

	schedules := []*SchedulesResponseItem{}
	if err := assertJSONBody(res, &schedules); err != nil {
		step.AddError(err)
	}

	return schedules, nil
}

type ScheduleResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Capacity     uint   `json:"capacity"`
	Reserved     uint   `json:"reserved"`
	Reservations []struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
		User   struct {
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
		} `json:"user"`
	} `json:"reservations"`
}

func ActionGetSchedule(ctx context.Context, step *isucandar.BenchmarkStep, id string, user *User) (*ScheduleResponse, error) {
	req, err := user.Agent.GET("/api/schedules/" + id)
	if err != nil {
		return nil, failure.NewError(ErrCritical, err)
	}

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
	}

	if err := assertContentType(res, "application/json"); err != nil {
		step.AddError(err)
	}

	schedule := &ScheduleResponse{}
	if err := assertJSONBody(res, schedule); err != nil {
		step.AddError(err)
	} else {
		if !user.Staff {
			for _, rev := range schedule.Reservations {
				if rev.User.Email != "" {
					step.AddError(failure.NewError(ErrSecurityIncident, fmt.Errorf("Leakage of email addresses at /schedules/%s", id)))
					break
				}
			}
		}
	}

	return schedule, nil
}

func ActionCreateReservation(ctx context.Context, step *isucandar.BenchmarkStep, schedule *Schedule, user *User) error {
	values := url.Values{}
	values.Add("schedule_id", schedule.ID)

	body := strings.NewReader(values.Encode())
	req, err := user.Agent.POST("/api/reservations", body)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := user.Agent.Do(ctx, req)
	if err != nil {
		return err
	}

	hasError := false
	if err := assertStatusCode(res, 200); err != nil {
		// 200 じゃなかったらなんか駄目だったんだな、という判定
		// step.AddError(err)
		hasError = true
	}

	if !hasError {
		user.mu.Lock()
		user.ReservedScheduleIDs = append(user.ReservedScheduleIDs, schedule.ID)
		user.mu.Unlock()

		schedule.Users.Add(user)

		step.AddScore(ScoreCreateReservation)
	}

	return nil
}
