package util

import (
	"crypto/rand"
	"log"
)

func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate random bytes:", err)
	}
	return b
}
