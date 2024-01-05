package random

import (
	"math/rand"
	"time"
)

// NewRandomString generates a random string with a given size.
// It uses a secure random number generator based on the current time.
// The string can contain uppercase letters, lowercase letters, and digits.
func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
