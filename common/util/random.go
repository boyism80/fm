// Package util provides common utility functions for the MapleStory private server.
// This file contains cryptographically secure random number generation utilities.
package util

import (
	"crypto/rand"
	"log"
)

// GenerateRandomBytes generates cryptographically secure random bytes.
// Used for session tokens, encryption keys, and other security-critical values.
//
// Flow: Called by various security components -> crypto/rand -> return bytes
// Usage: Session management, encryption, packet security
// Security: Uses crypto/rand for cryptographically secure randomness
//
// Parameters:
//   - size: Number of random bytes to generate (must be positive)
//
// Returns:
//   - []byte: Slice containing the generated random bytes
//
// Panics:
//   - If the system's random number generator fails (critical security failure)
func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate random bytes:", err)
	}
	return b
}
