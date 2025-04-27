package util

import (
	"fmt"
	"strings"
)

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
