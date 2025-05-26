package req

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/common/stream"
)

type Login struct {
	ID  string
	Pw  string
	Mac string
}

func (a *Login) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Login) Deserialize(reader *stream.StreamReader) error {
	id, err := reader.ReadStr16()
	if err != nil {
		return err
	}
	a.ID = id

	pw, err := reader.ReadStr16()
	if err != nil {
		return err
	}
	a.Pw = pw

	mac, err := reader.Read(6)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for i, b := range mac {
		if i > 0 {
			sb.WriteString("-")
		}
		sb.WriteString(fmt.Sprintf("%02X", b))
	}
	a.Mac = sb.String()
	return nil
}
