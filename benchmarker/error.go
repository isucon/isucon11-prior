package main

import (
	"context"
	"net"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
)

// Critical Errors
var (
	ErrCritical failure.StringCode = "critical"
)

func isCritical(err error) bool {
	// Prepare step でのエラーはすべて Critical の扱い
	return failure.IsCode(err, isucandar.ErrPrepare) ||
		failure.IsCode(err, ErrCritical)
}

var (
	ErrInvalidStatusCode  failure.StringCode = "status code"
	ErrInvalidContentType failure.StringCode = "content type"
	ErrInvalidJSON        failure.StringCode = "json"
	ErrMissmatch          failure.StringCode = "missmatch"
	ErrInvalidAsset       failure.StringCode = "asset"
	ErrInvalid            failure.StringCode = "invalid"
)

func isDeduction(err error) bool {
	return failure.IsCode(err, ErrInvalidStatusCode) ||
		failure.IsCode(err, ErrInvalidContentType) ||
		failure.IsCode(err, ErrInvalidJSON) ||
		failure.IsCode(err, ErrInvalidAsset) ||
		failure.IsCode(err, ErrMissmatch) ||
		failure.IsCode(err, ErrInvalid)
}

func isTimeout(err error) bool {
	var nerr net.Error
	if failure.As(err, &nerr) {
		if nerr.Timeout() || nerr.Temporary() {
			return true
		}
	}
	if failure.Is(err, context.DeadlineExceeded) {
		return true
	}
	return failure.IsCode(err, failure.TimeoutErrorCode)
}
