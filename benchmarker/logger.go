package main

import (
	"log"
	"os"
)

var (
	// 競技者に見せてもいい内容を書くロガー
	ContestantLogger *log.Logger
	// 運営だけが見れる内容を書くロガー
	AdminLogger *log.Logger
)

func init() {
	ContestantLogger = log.New(os.Stdout, "", log.Lmicroseconds)
	AdminLogger = log.New(os.Stderr, "", log.Lmicroseconds)
}

type Blackhole struct{}

func (b *Blackhole) Write(p []byte) (n int, err error) {
	return 0, nil
}
