package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type MultiChat struct {
	Mode    pconst.MultiChatMode
	Name    string
	Message string
}

func (p *MultiChat) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Mode))
	w.WriteStr16(p.Name)
	w.WriteStr16(p.Message)
	return nil
}

func (p *MultiChat) Deserialize(*stream.StreamReader) {}

func (p *MultiChat) Opcode() uint16 {
	return 0x5C
}
