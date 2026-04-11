package request

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/stream"
)

type Login struct {
	ID  string
	Pw  string
	Mac string
}

func (a *Login) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Login) Deserialize(reader *stream.StreamReader) {
	a.ID = reader.ReadStr16()
	a.Pw = reader.ReadStr16()

	mac := reader.Read(6)

	var sb strings.Builder
	for i, b := range mac {
		if i > 0 {
			sb.WriteString("-")
		}
		sb.WriteString(fmt.Sprintf("%02X", b))
	}
	a.Mac = sb.String()
}
