package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"

	"github.com/isucon/isucandar/agent"

	"github.com/isucon/isucandar"
)

type Scenario struct {
	// TODO: シナリオ実行に必要なフィールドを書く

	BaseURL     string // ベンチ対象 Web アプリの URL
	UseTLS      bool   // https で接続するかどうか
	NoLoad      bool   // Load(ベンチ負荷)を強要しない
	Parallelism int32  // リクエスト並列数

	// 競技者の実装言語
	Language string

	StaffUser *User
	Users     *Users
	Schedules *Schedules
}

func NewScenario() (*Scenario, error) {
	return &Scenario{
		UseTLS:    false,
		NoLoad:    false,
		Language:  "",
		Users:     newUsers(),
		Schedules: newSchedules(),
	}, nil
}

func (s *Scenario) Prepare(ctx context.Context, step *isucandar.BenchmarkStep) error {
	/*
		TODO: 負荷走行前の初期化部分をここに書く(ex: GET /initialize とか)
	*/
	initializer, err := agent.NewAgent(agent.WithBaseURL(s.BaseURL))
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	req, err := initializer.POST("/initialize", nil)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	res, err := initializer.Do(ctx, req)
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}

	assertInitialize(step, res)

	// スケジュール作ったりする管理ユーザー
	// Agent は initialize で使った奴使いまわしちゃおう
	s.StaffUser = &User{
		Email:    "isucon2021_prior@isucon.net",
		Nickname: "isucon",
		Staff:    true,
		Agent:    initializer,
	}

	if err := ActionLogin(ctx, step, s.StaffUser); err != nil {
		return err
	}

	return nil
}

func (s *Scenario) Load(parent context.Context, step *isucandar.BenchmarkStep) error {
	if s.NoLoad {
		return nil
	}

	ctx, cancel := context.WithTimeout(parent, 60*time.Second)
	defer cancel()

	userTimer, timerCancel := context.WithTimeout(ctx, 45*time.Second)
	defer timerCancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// すべてのスケジュールが予定で埋まっていたら、新しくスケジュールを作る
	scheduleWorker, err := worker.NewWorker(func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}

		wg.Add(1)
		defer wg.Done()

		err := BrowserAccess(ctx, step, s.StaffUser, "/")
		if err != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		schedules, err := ActionGetSchedules(ctx, step, s.StaffUser)
		if err != nil {
			step.AddError(err)
			return
		}

		hasCapacity := 0
		for _, s := range schedules {
			if s.Capacity > s.Reserved {
				hasCapacity++
			}
		}
		if hasCapacity >= 6 {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		schedule, err := ActionCreateSchedule(ctx, step, s)
		if err != nil {
			step.AddError(err)
		} else {
			s.Schedules.Add(schedule)
		}
	}, worker.WithInfinityLoop(), worker.WithMaxParallelism(s.Parallelism))
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}
	go func() {
		wg.Add(1)
		scheduleWorker.Process(ctx)
		wg.Done()
	}()

	userWorker, err := worker.NewWorker(func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		case <-userTimer.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}

		wg.Add(1)
		defer wg.Done()

		user, err := s.NewUser()
		if err != nil {
			step.AddError(failure.NewError(ErrCritical, err))
			return
		}

		if err := BrowserAccess(ctx, step, user, "/"); err != nil {
			step.AddError(err)
			return
		}
		if err := ActionSignup(ctx, step, user); err != nil {
			step.AddError(err)
			return
		}

		if user.ID == "" {
			return
		}

		if err := ActionLogin(ctx, step, user); err != nil {
			step.AddError(err)
			return
		}
		// ログイン失敗する人はここでおしまい
		if user.FailOnLogin {
			return
		}

		for !user.IsEnoughNeeds() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// リロードして
			if err := BrowserAccess(ctx, step, user, "/"); err != nil {
				step.AddError(err)
				continue
			}

			select {
			case <-ctx.Done():
				return
			default:
			}

			// スケジュール一覧を見て
			schedules, err := ActionGetSchedules(ctx, step, user)
			if err != nil {
				step.AddError(err)
				continue
			}

			for _, schedule := range schedules {
				if schedule.Capacity <= schedule.Reserved {
					continue
				}
				sschedule := s.Schedules.GetByID(schedule.ID)
				if sschedule == nil {
					step.AddError(failure.NewError(ErrInvalid, fmt.Errorf("存在しないはずのスケジュール ID です: %s", schedule.ID)))
					continue
				}

				if user.IsReserved(sschedule) {
					continue
				}

				select {
				case <-ctx.Done():
					return
				default:
				}

				// キャパが空いてて、取ってない予定なら抑えにかかる
				rschedule, err := ActionGetSchedule(ctx, step, schedule.ID, user)
				if err != nil {
					step.AddError(err)
					break
				}
				if rschedule.Capacity <= rschedule.Reserved {
					continue
				}

				select {
				case <-ctx.Done():
					return
				default:
				}

				if err := ActionCreateReservation(ctx, step, sschedule, user); err != nil {
					step.AddError(err)
					break
				}
			}
		}
	}, worker.WithInfinityLoop(), worker.WithMaxParallelism(s.Parallelism))
	if err != nil {
		return failure.NewError(ErrCritical, err)
	}
	userWorker.Process(ctx)

	wg.Done()
	wg.Wait()

	return nil
}

func (s *Scenario) Validation(parent context.Context, step *isucandar.BenchmarkStep) error {
	if s.NoLoad {
		return nil
	}

	// 10 秒待つ
	time.Sleep(10 * time.Second)

	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()

	/*
		- スケジュール数の一致
		- reservation の一致
	*/
	err := BrowserAccess(ctx, step, s.StaffUser, "/")
	if err != nil {
		return err
	}

	schedules, err := ActionGetAllSchedules(ctx, step, s.StaffUser)
	if err != nil {
		return err
	}

	if err := assertEqualInt(s.Schedules.Count(), len(schedules), "all-schedules.count"); err != nil {
		step.AddError(err)
	}

	scheduleWorker, err := worker.NewWorker(func(ctx context.Context, i int) {
		schedule := schedules[i]

		sschedule := s.Schedules.GetByID(schedule.ID)
		if sschedule == nil {
			step.AddError(failure.NewError(ErrInvalid, fmt.Errorf("unknown schedule id: %s", schedule.ID)))
			return
		}

		if err := assertEqualString(sschedule.Title, schedule.Title, "all-schedules.title"); err != nil {
			step.AddError(err)
		}

		if err := assertEqualUint(sschedule.Capacity, schedule.Capacity, "all-schedules.capacity"); err != nil {
			step.AddError(err)
		}

		if sschedule.Users.Count() > int(schedule.Reserved) {
			step.AddError(failure.NewError(ErrMissmatch, fmt.Errorf("schedule.reserved %d != %d", schedule.Reserved, sschedule.Users.Count())))
		}

		resp, err := ActionGetSchedule(ctx, step, schedule.ID, s.StaffUser)
		if err != nil {
			step.AddError(err)
			return
		}

		if err := assertEqualString(sschedule.Title, resp.Title, "schedule.title"); err != nil {
			step.AddError(err)
		}

		if err := assertEqualUint(sschedule.Capacity, resp.Capacity, "schedule.capacity"); err != nil {
			step.AddError(err)
		}

		if sschedule.Users.Count() > int(resp.Reserved) {
			step.AddError(failure.NewError(ErrMissmatch, fmt.Errorf("schedule.reserved %d != %d", resp.Reserved, sschedule.Users.Count())))
		}

		allowUnknownUsersCount := 0
		if err := assertEqualInt(sschedule.Users.Count(), len(resp.Reservations), "schedule.reservations.count"); err != nil {
			step.AddError(err)
			allowUnknownUsersCount = len(resp.Reservations) - sschedule.Users.Count()
		}

		if len(resp.Reservations) > int(sschedule.Capacity) {
			step.AddError(failure.NewError(ErrInvalid, fmt.Errorf("overbooking at %s", sschedule.ID)))
		}

		revMap := map[string]string{}
		for _, reservation := range resp.Reservations {
			suser := sschedule.Users.GetByID(reservation.UserID)
			if suser == nil {
				if allowUnknownUsersCount <= 0 {
					step.AddError(failure.NewError(ErrInvalid, fmt.Errorf("unknown user on reservations: %s", reservation.UserID)))
				} else {
					allowUnknownUsersCount--
				}
				continue
			}

			if revID, ok := revMap[suser.ID]; ok {
				step.AddError(failure.NewError(ErrInvalid, fmt.Errorf("duplication reservation on schedule(id: %s) / reservation(id: %s)", sschedule.ID, revID)))
			}
			revMap[suser.ID] = reservation.ID

			if err := assertEqualString(suser.Email, reservation.User.Email, "reservation.user.email"); err != nil {
				step.AddError(err)
			}

			if err := assertEqualString(suser.Nickname, reservation.User.Nickname, "reservation.user.nickname"); err != nil {
				step.AddError(err)
			}
		}
	}, worker.WithLoopCount(int32(len(schedules))))
	if err != nil {
		return err
	}
	scheduleWorker.Process(ctx)

	return nil
}

func (s *Scenario) NewUser() (*User, error) {
	a, err := agent.NewAgent(agent.WithBaseURL(s.BaseURL))
	if err != nil {
		return nil, failure.NewError(ErrCritical, err)
	}

	user := newUser()
	user.Agent = a

	return user, nil
}

func (s *Scenario) NewSchedule() (*Schedule, error) {
	return newSchedule(), nil
}
