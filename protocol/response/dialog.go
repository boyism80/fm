package response

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/services/game/constant"
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

type DialogStyle struct {
	NPC    uint32
	Text   string
	Styles []uint32
}

func (p *DialogInput) Opcode() uint16 {
	return 0xE5
}

func (p *DialogStyle) Opcode() uint16 {
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
	writer.WriteU8(uint8(constant.DialogTypeYesNo))
	writer.WriteStr16(p.Text)
	writer.WriteBoolean(p.Prev)
	writer.WriteBoolean(p.Next)
	return nil
}

func (p *DialogList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DialogTypeList))

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
		writer.WriteU8(uint8(constant.DialogTypeAcceptEscape))
	} else {
		writer.WriteU8(uint8(constant.DialogTypeAccept))
	}
	writer.WriteStr16(p.Text)
	return nil
}

func (p *DialogInput) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DialogTypeInput))
	writer.WriteStr16(p.Text)
	writer.WriteU32(0)
	writer.WriteU32(0)
	return nil
}

func (p *DialogStyle) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.WriteU32(p.NPC)
	writer.WriteU8(uint8(constant.DialogTypeStyle))
	writer.WriteStr16(p.Text)
	writer.WriteU8(uint8(len(p.Styles)))
	for _, style := range p.Styles {
		writer.WriteU32(style)
	}
	return nil
}

func (p *Dialog) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	p.Type = constant.DialogType(reader.ReadU8())
	p.Text = reader.ReadStr16()
	p.Prev = reader.ReadBool()
	p.Next = reader.ReadBool()
}

func (p *DialogYesNo) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	reader.ReadU8()
	p.Text = reader.ReadStr16()
	p.Prev = reader.ReadBool()
	p.Next = reader.ReadBool()
}

func (p *DialogList) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	reader.ReadU8()
	parts := strings.Split(reader.ReadStr16(), "\r\n#b#L")
	p.Text = parts[0]
	p.Selections = nil
	for _, part := range parts[1:] {
		_, selection, _ := strings.Cut(part, "# ")
		p.Selections = append(p.Selections, strings.TrimSuffix(selection, "#l"))
	}
}

func (p *DialogAccept) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	p.EnableEscape = constant.DialogType(reader.ReadU8()) == constant.DialogTypeAcceptEscape
	p.Text = reader.ReadStr16()
}

func (p *DialogInput) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	reader.ReadU8()
	p.Text = reader.ReadStr16()
	reader.ReadU32()
	reader.ReadU32()
}

func (p *DialogStyle) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NPC = reader.ReadU32()
	reader.ReadU8()
	p.Text = reader.ReadStr16()
	n := int(reader.ReadU8())
	p.Styles = make([]uint32, 0, n)
	for i := 0; i < n; i++ {
		p.Styles = append(p.Styles, reader.ReadU32())
	}
}
