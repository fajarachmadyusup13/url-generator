package utils

import (
	"math/rand"
	"time"
)

func GenerateID() int64 {
	return time.Now().UnixNano() + int64(rand.Intn(10000))
}

func GenerateString() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, 10)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}
