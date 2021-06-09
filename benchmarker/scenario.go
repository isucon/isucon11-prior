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
}

func NewScenario() (*Scenario, error) {
	return &Scenario{
		// TODO: シナリオを初期化する
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

	// ここで res を検証する
	assertInitialize(step, res)

	return nil
}

func (s *Scenario) Load(parent context.Context, step *isucandar.BenchmarkStep) error {
	if s.NoLoad {
		return nil
	}

	/*
		TODO: 実際の負荷走行シナリオ
	*/

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
