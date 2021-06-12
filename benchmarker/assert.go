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
		return failure.NewError(ErrInvalidStatusCode, fmt.Errorf("invalid status code: %d (expected: %d) at %s", res.StatusCode, code, res.Request.URL.Path))
	}
	return nil
}

func assertContentType(res *http.Response, contentType string) error {
	actual := res.Header.Get("Content-Type")
	if !strings.HasPrefix(actual, contentType) {
		return failure.NewError(ErrInvalidContentType, fmt.Errorf("invalid content type: %s (expected: %s) at %s", actual, contentType, res.Request.URL.Path))
	}
	return nil
}

func assertJSONBody(res *http.Response, body interface{}) error {
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	if err := decoder.Decode(body); err != nil {
		return failure.NewError(ErrInvalidJSON, fmt.Errorf("invalid JSON at %s", res.Request.URL.Path))
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

func assertEqualString(expected, actual string, tag string) error {
	if expected != actual {
		return failure.NewError(ErrMissmatch, fmt.Errorf("%s %s != expected %s", tag, actual, expected))
	}
	return nil
}

func assertEqualInt(expected, actual int, tag string) error {
	if expected != actual {
		return failure.NewError(ErrMissmatch, fmt.Errorf("%s %d != expected %d", tag, actual, expected))
	}
	return nil
}

func assertEqualUint(expected, actual uint, tag string) error {
	if expected != actual {
		return failure.NewError(ErrMissmatch, fmt.Errorf("%s %d != expected %d", tag, actual, expected))
	}
	return nil
}
