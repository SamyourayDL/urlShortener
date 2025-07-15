package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandomString(length int) string {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	b.Grow(length)
	for i := 0; i < length; i++ {
		b.WriteRune(chars[gen.Intn(len(chars))])
	}

	return b.String()
}
