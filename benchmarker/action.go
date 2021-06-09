package main

import (
	"context"
	"net/url"
	"strings"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
)

func ActionSignup(step *isucandar.BenchmarkStep, a *agent.Agent) (*User, error) {
	return nil, nil
}

// Action がエラーを返す → Action の失敗
// Action がエラーを返さない → Action としては成功。シナリオとしてはどうかわからない
func ActionLogin(ctx context.Context, step *isucandar.BenchmarkStep, u *User) error {
	values := url.Values{}
	values.Add("email", u.Email)

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

	if err := assertStatusCode(res, 200); err != nil {
		step.AddError(err)
		return nil
	}

	return nil
}

func ActionCreateSchedule(step *isucandar.BenchmarkStep, a *agent.Agent) {}
