package byte_trie_test

import (
	"math/rand"
	"time"
)

var defaultChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ 0123456789")

func randomString(maxLength int, chars ...rune) string {
	if len(chars) == 0 {
		chars = defaultChars
	}

	length := rand.Intn(maxLength) + 1
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func randomStrings(maxLength, count int, chars ...rune) []string {
	if len(chars) == 0 {
		chars = defaultChars
	}

	rand.Seed(time.Now().UnixNano())

	ss := make([]string, count)
	for i := 0; i < count; i++ {
		ss[i] = randomString(maxLength, chars...)
	}

	return ss
}
