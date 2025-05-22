package utils

import (
	"math/rand/v2"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for range n {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}

	return sb.String()
}

func RandomPassword() string {
	parts := []string{RandomString(4), RandomString(4), RandomString(4)}
	return strings.Join(parts, "-")
}
