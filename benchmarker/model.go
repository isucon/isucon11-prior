package main

import (
	"sync"
	"time"

	"github.com/isucon/isucandar/agent"
)

// DB のスキーマと合わせつつ、ベンチマーカーが検証に利用するためのデータモデル

type User struct {
	ID        string
	Email     string
	Nickname  string
	Staff     bool
	CreatedAt time.Time

	// 1ユーザーごとに Cookie を持つので 1 ユーザーごとに Agent を専有したほうがいい
	// Agent 1つ と UserAgent (ブラウザ) が1:1になるイメージ
	Agent *agent.Agent
}

func newUser() *User {
	return &User{
		ID:        "",
		Email:     randomEmail(),
		Nickname:  randomNickname(),
		Staff:     false,
		CreatedAt: time.Unix(0, 0),
		Agent:     nil,
	}
}

type Users struct {
	mu    sync.Mutex
	slice []*User
}

func newUsers() *Users {
	return &Users{
		mu:    sync.Mutex{},
		slice: []*User{},
	}
}

func (a *Users) Add(u *User) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.slice = append(a.slice, u)
}

type Schedule struct {
	ID        string
	Title     string
	Capacity  uint
	CreatedAt time.Time
}

type Reservation struct {
	ID        string
	Schedule  *Schedule
	User      *User
	CreatedAt time.Time
}
