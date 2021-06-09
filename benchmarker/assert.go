package main

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
)

// 検証用コード群

func assertInitialize(step *isucandar.BenchmarkStep, res *http.Response) {
	err := assertStatusCode(res, 200)
	if err != nil {
		step.AddError(failure.NewError(ErrCritical, err))
	}
}

func assertStatusCode(res *http.Response, code int) error {
	if res.StatusCode != code {
		return failure.NewError(ErrInvalidStatusCode, fmt.Errorf("Invalid status code: %d (expected: %d)", res.StatusCode, code))
	}
	return nil
}
