package main

import "github.com/isucon/isucandar/failure"

// Critical Errors
var (
	ErrCritical failure.StringCode = "CRITICAL"
)

func isCritical(err error) bool {
	return failure.IsCode(err, ErrCritical)
}

var (
	ErrInvalidStatusCode failure.StringCode = "INVALID STATUS CODE"
)
