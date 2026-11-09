package wz

import (
	"encoding/xml"
	"os"
	"strings"
)

const (
	nameMinBytes = 3
	nameMaxBytes = 12
)

func loadForbiddenNames(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	words := make([]string, 0, len(root.Strings))
	// Curse.img values are "word,replacement"; the word before the comma is forbidden.
	for _, s := range root.Strings {
		if word, _, _ := strings.Cut(s.Value, ","); word != "" {
			words = append(words, strings.ToLower(word))
		}
	}
	return words, nil
}

// ValidName is the client's rule for names: letters, digits and Hangul syllables
// only, 3-12 bytes in CP949 (a Hangul syllable takes 2), and no forbidden word.
func (r *Resources) ValidName(name string) bool {
	bytes := 0
	for _, c := range name {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
			bytes++
		case c >= '가' && c <= '힣':
			bytes += 2
		default:
			return false
		}
	}
	if bytes < nameMinBytes || bytes > nameMaxBytes {
		return false
	}
	lower := strings.ToLower(name)
	for _, word := range r.ForbiddenNames {
		if strings.Contains(lower, word) {
			return false
		}
	}
	return true
}
