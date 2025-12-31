// Package util provides common utility functions for the MapleStory private server.
// This file contains string manipulation utilities.
package util

import (
	"fmt"
	"strings"
)

// ToHexString converts a byte slice to a hex string representation.
// Each byte is formatted as uppercase hex with spaces between bytes.
func ToHexString(bytes []byte) string {
	var hexed strings.Builder
	for i := 0; i < len(bytes); i++ {
		hexed.WriteString(fmt.Sprintf("%02X ", bytes[i]))
	}
	result := hexed.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result
}
