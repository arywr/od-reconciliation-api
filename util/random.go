package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates random string with length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomType generates a random type name
func RandomType() string {
	types := []string{"cash-in", "cash-out", "refund", "additional", "other"}
	n := len(types)
	return types[rand.Intn(n)]
}

func RandomDescription() string {
	return RandomString(25)
}
