package main

import "fmt"

// ダミーデータの生成方法は最初のうちは凝らない方が楽です

var randomEmailCount int = 0

func randomEmail() string {
	randomEmailCount++
	return fmt.Sprintf("isucon-%d@example.com", randomEmailCount)
}

var randomNicknameCount int = 0

func randomNickname() string {
	randomNicknameCount++
	return fmt.Sprintf("isucon-%d", randomNicknameCount)
}
