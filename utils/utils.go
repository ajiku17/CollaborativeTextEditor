package utils

import "math/rand"

func SetRandSeed(seed int64) {
	rand.Seed(seed)
}

func RandBetween(low, high int) int {
	return rand.Intn(high-low) + low
}

func Min(x, y int) int {
	if x > y {
		return y
	}

	return x
}

func Max(x, y int) int {
	if x < y {
		return y
	}

	return x
}

func Reverse(s string) string {
	result := ""
	for _, ch := range s {
		result = string(ch) + result
	}
	return result
}
