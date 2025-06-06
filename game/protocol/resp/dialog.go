package resp

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type Dialog struct {
	NPC  uint32
	Type constant.DialogType
	Text string
	Prev bool
	Next bool
}

type DialogYesNo struct {
	NPC  uint32
	Text string
	Prev bool
	Next bool
}

type DialogList struct {
	NPC        uint32
	Text       string
	Selections []string
}

type DialogAccept struct {
	NPC          uint32
	Text         string
	EnableEscape bool
}

type DialogInput struct {
	NPC  uint32
	Text string
}

func (p *Dialog) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xE5)
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(p.Type))
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogYesNo) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xE5)
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DIALOG_TYPE_YES_NO))
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xE5)
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
	writer.WriteU16(0xE5)
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
	writer.WriteU16(0xE5)
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DIALOG_TYPE_INPUT))
	writer.WriteStr16(p.Text)
	writer.WriteU32(0)
	writer.WriteU32(0)
	return nil
}

func (p *Dialog) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *DialogYesNo) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *DialogList) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *DialogAccept) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *DialogInput) Deserialize(reader *stream.StreamReader) error {
	return nil
}
