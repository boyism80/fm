package request

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/boyism80/fm/stream"
)

type Login struct {
	ID  string
	Pw  string
	Mac string
}

func (*Login) Opcode() byte { return 0x01 }

func (a *Login) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.ID)
	writer.WriteStr16(a.Pw)
	parts := strings.Split(a.Mac, "-")
	if len(parts) != 6 {
		return fmt.Errorf("login mac %q", a.Mac)
	}
	mac := make([]byte, 6)
	for i, part := range parts {
		v, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return fmt.Errorf("login mac %q", a.Mac)
		}
		mac[i] = byte(v)
	}
	writer.Write(mac)
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
