package server

import (
	"strings"
	"unicode"

	"golang.org/x/text/encoding/korean"
)

type socialText string

// valid measures the limit in EUC-KR bytes, the client's text encoding.
func (t socialText) valid(limit int) bool {
	if strings.TrimSpace(string(t)) == "" {
		return false
	}

	for _, r := range t {
		if unicode.IsControl(r) {
			return false
		}
	}

	encoded, err := korean.EUCKR.NewEncoder().Bytes([]byte(t))
	return err == nil && len(encoded) <= limit
}
