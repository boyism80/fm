package response

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type Dialog struct {
	NPC  uint32
	Type constant.DialogType
	Text string
	Prev bool
	Next bool
}

func (p *Dialog) Opcode() uint16 {
	return 0xE5
}

type DialogYesNo struct {
	NPC  uint32
	Text string
	Prev bool
	Next bool
}

func (p *DialogYesNo) Opcode() uint16 {
	return 0xE5
}

type DialogList struct {
	NPC        uint32
	Text       string
	Selections []string
}

func (p *DialogList) Opcode() uint16 {
	return 0xE5
}

type DialogAccept struct {
	NPC          uint32
	Text         string
	EnableEscape bool
}

func (p *DialogAccept) Opcode() uint16 {
	return 0xE5
}

type DialogInput struct {
	NPC  uint32
	Text string
}

func (p *DialogInput) Opcode() uint16 {
	return 0xE5
}

func (p *Dialog) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(p.Type))
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogYesNo) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DIALOG_TYPE_YES_NO))
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DIALOG_TYPE_LIST))

	var builder strings.Builder
	builder.WriteString(p.Text)
	for i, selection := range p.Selections {
		builder.WriteString("\r\n")
		builder.WriteString(fmt.Sprintf("#b#L%d# %s#l", i, selection))
	}
	writer.WriteStr16(builder.String())
	return nil
}

func (p *DialogAccept) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	if p.EnableEscape {
		writer.WriteU8(uint8(constant.DIALOG_TYPE_ACCEPT_ESCAPE))
	} else {
		writer.WriteU8(uint8(constant.DIALOG_TYPE_ACCEPT))
	}
	writer.WriteStr16(p.Text)
	return nil
}

func (p *DialogInput) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DIALOG_TYPE_INPUT))
	writer.WriteStr16(p.Text)
	writer.WriteU32(0)
	writer.WriteU32(0)
	return nil
}

func (p *Dialog) Deserialize(reader *stream.StreamReader) {
}

func (p *DialogYesNo) Deserialize(reader *stream.StreamReader) {
}

func (p *DialogList) Deserialize(reader *stream.StreamReader) {
}

func (p *DialogAccept) Deserialize(reader *stream.StreamReader) {
}

func (p *DialogInput) Deserialize(reader *stream.StreamReader) {
}
