package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/manveru/faker"
)

// ダミーデータの生成方法は最初のうちは凝らない方が楽です

var (
	emailLock  = sync.Mutex{}
	emailFaker *faker.Faker
	nameLock   = sync.Mutex{}
	nameFaker  *faker.Faker
)

func init() {
	var err error
	emailFaker, err = faker.New("en")
	if err != nil {
		panic(err)
	}
	nameFaker, err = faker.New("en")
	if err != nil {
		panic(err)
	}
}

// 一定確率で true
func percentage(numerator int, denominator int) bool {
	return rand.Intn(denominator) <= numerator
}

var randomEmailCount int64 = 0

// インクリメントで race したので直す
func randomEmail() string {
	emailLock.Lock()
	defer emailLock.Unlock()
	cnt := atomic.AddInt64(&randomEmailCount, 1)
	return fmt.Sprintf("%d-%s", cnt, emailFaker.Email())
}

func randomNickname() string {
	nameLock.Lock()
	defer nameLock.Unlock()
	return nameFaker.Name()
}

func randomTitle() string {
	return fmt.Sprintf("%s %s", randomDate(), randomCity())
}

func randomCapacity() int {
	return 30 + rand.Intn(100)
}

var randomDateCount int32 = 0

func randomDate() string {
	cnt := atomic.AddInt32(&randomDateCount, 1)
	date := time.Now().AddDate(0, 0, int(cnt))
	return date.Format("2006-01-02")
}
