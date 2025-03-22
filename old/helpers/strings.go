package helpers

import (
	"math/rand"
)

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		letter := letters[rand.Int()%len(letters)]
		result[i] = letter
	}
	return string(result)
}
