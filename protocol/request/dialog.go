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
	case constant.DIALOG_TYPE_DEFAULT:
	case constant.DIALOG_TYPE_YES_NO:

	case constant.DIALOG_TYPE_LIST:
		p.Selected = reader.ReadU32()

	case constant.DIALOG_TYPE_INPUT:
		p.Text = reader.ReadStr16()

	case constant.DIALOG_TYPE_ACCEPT_ESCAPE:
	case constant.DIALOG_TYPE_ACCEPT:
	default:
		panic(fmt.Errorf("invalid dialog type: %d", p.DialogType))
	}

}
