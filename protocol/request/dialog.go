package request

import (
	"fmt"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type Dialog struct {
	DialogType constant.DialogType
	Selected   uint32
	Next       bool
	Text       string
}

func (*Dialog) Opcode() byte { return 0x2B }

func (p *Dialog) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Dialog) Deserialize(reader *stream.StreamReader) {
	p.DialogType = constant.DialogType(reader.ReadU8())
	p.Next = reader.ReadU8() == 1

	switch p.DialogType {
	case constant.DialogTypeDefault:
	case constant.DialogTypeYesNo:

	case constant.DialogTypeList:
		p.Selected = reader.ReadU32()

	case constant.DialogTypeInput:
		p.Text = reader.ReadStr16()

	case constant.DialogTypeAcceptEscape:
	case constant.DialogTypeAccept:
	default:
		panic(fmt.Errorf("invalid dialog type: %d", p.DialogType))
	}

}
