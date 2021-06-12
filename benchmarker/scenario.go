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
		if err := ActionSignup(ctx, step, user); err != nil || user.ID == "" {
			step.AddError(err)
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

				// キャパが空いてて、取ってない予定なら抑えにかかる
				rschedule, err := ActionGetSchedule(ctx, step, schedule.ID, user)
				if err != nil {
					step.AddError(err)
					break
				}
				if rschedule.Capacity <= rschedule.Reserved {
					continue
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

func (s *Scenario) Validation(ctx context.Context, step *isucandar.BenchmarkStep) error {
	if s.NoLoad {
		return nil
	}

	/*
		TODO: 負荷走行後のデータ検証シナリオ
	*/

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
