package resp

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/common/stream"
)

type DialogType uint8

const (
	DIALOG_TYPE_DEFAULT DialogType = 0
	DIALOG_TYPE_YES_NO  DialogType = 1
	DIALOG_TYPE_LIST    DialogType = 4
	DIALOG_TYPE_ACCEPT  DialogType = 11
)

type Dialog struct {
	NPC  uint32
	Type DialogType
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
	writer.WriteU8(1)
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xE5)
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(4)

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
		writer.WriteU8(11)
	} else {
		writer.WriteU8(12)
	}
	writer.WriteStr16(p.Text)
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
