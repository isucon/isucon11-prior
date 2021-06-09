package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func assertContentType(res *http.Response, contentType string) error {
	actual := res.Header.Get("Content-Type")
	if !strings.HasPrefix(actual, contentType) {
		return failure.NewError(ErrInvalidContentType, fmt.Errorf("Invalid content type: %s (expected: %s)", actual, contentType))
	}
	return nil
}

func assertJSONBody(res *http.Response, body interface{}) error {
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	if err := decoder.Decode(body); err != nil {
		return failure.NewError(ErrInvalidJSON, fmt.Errorf("Invalid JSON"))
	}
	return nil
}

func assertChecksum(res *http.Response) error {
	defer res.Body.Close()

	path := res.Request.URL.Path
	expected := resoucesHash[path]
	if expected == "" {
		return nil
	}

	hash := md5.New()
	if _, err := io.Copy(hash, res.Body); err != nil {
		return failure.NewError(ErrCritical, err)
	}
	actual := fmt.Sprintf("%x", hash.Sum(nil))

	if expected != actual {
		return failure.NewError(ErrInvalidAsset, fmt.Errorf("invalid MD5: %s %s != expected %s", path, actual, expected))
	}
	return nil
}

func assertEqualString(expected, actual string) error {
	if expected != actual {
		return failure.NewError(ErrMissmatch, fmt.Errorf("missmatch: %s != expected %s", actual, expected))
	}
	return nil
}

func assertEqualUint(expected, actual uint) error {
	if expected != actual {
		return failure.NewError(ErrMissmatch, fmt.Errorf("missmatch: %d != expected %d", actual, expected))
	}
	return nil
}
