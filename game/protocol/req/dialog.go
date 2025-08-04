package req

import (
	"fmt"

	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/constant"
)

type Dialog struct {
	DialogType constant.DialogType
	Selected   uint32
	Next       bool
	Text       string
}

func (p *Dialog) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Dialog) Deserialize(reader *stream.StreamReader) error {
	dialogType, _ := reader.ReadU8()
	p.DialogType = constant.DialogType(dialogType)
	action, _ := reader.ReadU8()
	p.Next = action == 1
	if !p.Next {
		return nil
	}

	switch p.DialogType {
	case constant.DIALOG_TYPE_DEFAULT:
	case constant.DIALOG_TYPE_YES_NO:

	case constant.DIALOG_TYPE_LIST:
		p.Selected, _ = reader.ReadU32()

	case constant.DIALOG_TYPE_INPUT:
		p.Text, _ = reader.ReadStr16()

	case constant.DIALOG_TYPE_ACCEPT_ESCAPE:
	case constant.DIALOG_TYPE_ACCEPT:
	default:
		return fmt.Errorf("invalid dialog type: %d", p.DialogType)
	}

	return nil
}
