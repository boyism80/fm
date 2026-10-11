package entity

import (
	"strings"
	"unicode"

	"golang.org/x/text/encoding/korean"
)

type ClientText string

// Valid measures the limit in EUC-KR bytes, the client's text encoding.
func (t ClientText) Valid(limit int) bool {
	return strings.TrimSpace(string(t)) != "" && t.fits(limit)
}

// The client joins every input line with "\n", leaving empty lines in place.
func (t ClientText) ValidLines(count int, limit int) bool {
	lines := strings.Split(string(t), "\n")
	if len(lines) != count || strings.TrimSpace(string(t)) == "" {
		return false
	}
	for _, line := range lines {
		if !ClientText(line).fits(limit) {
			return false
		}
	}
	return true
}

func (t ClientText) fits(limit int) bool {
	for _, r := range t {
		if unicode.IsControl(r) {
			return false
		}
	}

	encoded, err := korean.EUCKR.NewEncoder().Bytes([]byte(t))
	return err == nil && len(encoded) <= limit
}
