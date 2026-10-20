package wz

import (
	"encoding/xml"
	"os"
	"strconv"
)

// OXQuestion is one question of Etc.wz/OXQuiz.img. The client shows the text itself from the set and id.
type OXQuestion struct {
	Set uint8
	ID  uint16
	O   bool
}

func loadOXQuiz(path string) ([]OXQuestion, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	var questions []OXQuestion
	for _, set := range root.Children {
		setID, err := strconv.Atoi(set.Name)
		if err != nil || setID <= 0 || setID > 255 {
			continue
		}
		for _, q := range set.Children {
			id, err := strconv.Atoi(q.Name)
			if err != nil || id <= 0 || id > 65535 {
				continue
			}
			questions = append(questions, OXQuestion{Set: uint8(setID), ID: uint16(id), O: q.Int("a", 0) == 1})
		}
	}
	return questions, nil
}
