package response

import (
	"github.com/boyism80/fm/stream"
)

// ShowSelfSafetyCharmEffect tells the client a safety charm kept the EXP a death would have taken.
type ShowSelfSafetyCharmEffect struct {
	Left     uint8
	DaysLeft uint8
}

func (p *ShowSelfSafetyCharmEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfSafetyCharmEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0x06)
	writer.WriteU8(1)
	writer.WriteU8(p.Left)
	writer.WriteU8(p.DaysLeft)
	return nil
}

func (p *ShowSelfSafetyCharmEffect) Deserialize(reader *stream.StreamReader) {
}
