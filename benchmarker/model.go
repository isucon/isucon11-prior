package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/isucon/isucandar/agent"
)

// DB のスキーマと合わせつつ、ベンチマーカーが検証に利用するためのデータモデル

type User struct {
	mu sync.Mutex

	ID        string
	Email     string
	Nickname  string
	Staff     bool
	CreatedAt time.Time

	// 1ユーザーごとに Cookie を持つので 1 ユーザーごとに Agent を専有したほうがいい
	// Agent 1つ と UserAgent (ブラウザ) が1:1になるイメージ
	Agent *agent.Agent

	FailOnSignup        bool
	FailOnLogin         bool
	Needs               int
	ReservedScheduleIDs []string
}

func newUser() *User {
	return &User{
		mu:                  sync.Mutex{},
		ID:                  "",
		Email:               randomEmail(),
		Nickname:            randomNickname(),
		Staff:               false,
		CreatedAt:           time.Unix(0, 0),
		Agent:               nil,
		FailOnSignup:        percentage(1, 100),
		FailOnLogin:         percentage(1, 20),
		Needs:               rand.Intn(3) + 3,
		ReservedScheduleIDs: []string{},
	}
}

func (u *User) IsReserved(s *Schedule) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, id := range u.ReservedScheduleIDs {
		if id == s.ID {
			return true
		}
	}
	return false
}

func (u *User) IsEnoughNeeds() bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	return len(u.ReservedScheduleIDs) >= u.Needs
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

func (a *Users) Get(idx int) *User {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.slice[idx]
}

func (a *Users) Count() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.slice)
}

type Schedule struct {
	ID        string
	Title     string
	Capacity  uint
	CreatedAt time.Time
}

func newSchedule() *Schedule {
	return &Schedule{
		ID:        "",
		Title:     randomTitle(),
		Capacity:  uint(randomCapacity()),
		CreatedAt: time.Unix(0, 0),
	}
}

type Schedules struct {
	mu    sync.Mutex
	slice []*Schedule
}

func newSchedules() *Schedules {
	return &Schedules{
		mu:    sync.Mutex{},
		slice: []*Schedule{},
	}
}

func (a *Schedules) Add(u *Schedule) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.slice = append(a.slice, u)
}

func (a *Schedules) Get(idx int) *Schedule {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.slice[idx]
}

func (a *Schedules) GetByID(id string) *Schedule {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, s := range a.slice {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (a *Schedules) Count() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.slice)
}

type Reservation struct {
	ID        string
	Schedule  *Schedule
	User      *User
	CreatedAt time.Time
}
