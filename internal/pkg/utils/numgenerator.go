package utils

import (
	"math/rand"
	"time"
)

func GenRandomNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
