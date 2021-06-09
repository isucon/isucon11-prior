package main

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"
)

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

	res, err := u.Agent.Do(ctx, req)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
		return nil
	}

	step.AddScore(ScoreSignup)

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
	}, worker.WithMaxParallelism(50), worker.WithInfinityLoop())
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

	res, err := u.Agent.Do(ctx, req)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	expectedStatusCode := 200
	if u.FailOnLogin {
		expectedStatusCode = 403
	}
	if err := assertStatusCode(res, expectedStatusCode); err != nil {
		step.AddError(err)
		return nil
	}

	step.AddScore(ScoreLogin)

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
	}, worker.WithMaxParallelism(100), worker.WithLoopCount(int32(usersCount)))
	if err != nil {
		return err
	}

	w.Process(ctx)

	wg.Done()
	wg.Wait()

	return nil
}

func ActionCreateSchedule(step *isucandar.BenchmarkStep, a *agent.Agent) {}
