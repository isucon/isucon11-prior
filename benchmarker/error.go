package main

import (
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
)

// Critical Errors
var (
	ErrCritical failure.StringCode = "CRITICAL"
)

func isCritical(err error) bool {
	// Prepare step でのエラーはすべて Critical の扱い
	return failure.IsCode(err, isucandar.ErrPrepare) || failure.IsCode(err, ErrCritical)
}

var (
	ErrInvalidStatusCode  failure.StringCode = "INVALID STATUS CODE"
	ErrInvalidContentType failure.StringCode = "INVALID CONTENT TYPE"
	ErrInvalidJSON        failure.StringCode = "INVALID JSON"
	ErrMissmatch          failure.StringCode = "MISSMATCH"
)
