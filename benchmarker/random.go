package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// ダミーデータの生成方法は最初のうちは凝らない方が楽です

// 一定確率で true
func percentage(numerator int, denominator int) bool {
	return rand.Intn(denominator) <= numerator
}

var randomEmailCount int64 = 0

// インクリメントで race したので直す
func randomEmail() string {
	cnt := atomic.AddInt64(&randomEmailCount, 1)
	return fmt.Sprintf("isucon-%d@example.com", cnt)
}

var randomNicknameCount int64 = 0

func randomNickname() string {
	cnt := atomic.AddInt64(&randomNicknameCount, 1)
	return fmt.Sprintf("isucon-%d", cnt)
}

func randomTitle() string {
	return fmt.Sprintf("%s %s", randomDate(), randomCity())
}

func randomCapacity() int {
	return 1 // 30 + rand.Intn(100)
}

var randomDateCount int32 = 0

func randomDate() string {
	cnt := atomic.AddInt32(&randomDateCount, 1)
	date := time.Now().AddDate(0, 0, int(cnt))
	return date.Format("2006-01-02")
}
