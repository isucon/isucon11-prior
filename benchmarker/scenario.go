package main

import (
	"context"

	"github.com/isucon/isucandar/failure"

	"github.com/isucon/isucandar/agent"

	"github.com/isucon/isucandar"
)

type Scenario struct {
	// TODO: シナリオ実行に必要なフィールドを書く

	BaseURL string // ベンチ対象 Web アプリの URL
	UseTLS  bool   // https で接続するかどうか
	NoLoad  bool   // Load(ベンチ負荷)を強要しない

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

	/*
		TODO: 実際の負荷走行シナリオ
	*/

	if err := ActionSignups(parent, step, s); err != nil {
		return err
	}

	if err := ActionLogins(parent, step, s); err != nil {
		return err
	}

	if err := ActionCreateSchedules(parent, step, s); err != nil {
		return err
	}

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
