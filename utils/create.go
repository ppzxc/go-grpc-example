package utils

import "crypto/rand"

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return b
}
