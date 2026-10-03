package conn

import (
	"fmt"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

var Responses = []any{
	&response.LoginFailed{},
	&response.Authenticate{},
	&response.ServerList{},
	&response.EndOfServerList{},
	&response.CharacterList{},
	&response.Transfer{},
	&response.CheckName{},
	&response.CreateCharacter{},
	&response.Login{},
	&response.Warp{},
	&response.KeyMap{},
	&response.Notice{},
	&response.SpawnNpc{},
	&response.RemoveNpc{},
	&response.Dialog{},
	&response.DialogYesNo{},
	&response.DialogInput{},
	&response.DialogList{},
	&response.DialogStyle{},
	&response.DialogAccept{},
}

func Decode(opcode uint16, body []byte) (any, error) {
	reader := stream.NewStreamReader(&body, stream.LittleEndian)
	switch opcode {
	case (&response.LoginFailed{}).Opcode():
		if len(body) == 0 {
			return nil, fmt.Errorf("empty login opcode")
		}
		if body[0] == 0 {
			pkt := &response.Authenticate{}
			pkt.Deserialize(reader)
			return pkt, nil
		}
		pkt := &response.LoginFailed{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.ServerList{}).Opcode():
		if len(body) > 0 && body[0] == 0xFF {
			pkt := &response.EndOfServerList{}
			pkt.Deserialize(reader)
			return pkt, nil
		}
		pkt := &response.ServerList{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.CharacterList{}).Opcode():
		pkt := &response.CharacterList{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.Transfer{}).Opcode():
		pkt := &response.Transfer{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.CheckName{}).Opcode():
		pkt := &response.CheckName{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.CreateCharacter{}).Opcode():
		pkt := &response.CreateCharacter{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.Ping{}).Opcode():
		return &response.Ping{}, nil
	case (&response.Login{}).Opcode():
		if len(body) > 4 && body[4] == 2 {
			pkt := &response.Warp{}
			pkt.Deserialize(reader)
			return pkt, nil
		}
		pkt := &response.Login{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.KeyMap{}).Opcode():
		pkt := &response.KeyMap{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.Notice{}).Opcode():
		pkt := &response.Notice{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.SpawnNpc{}).Opcode():
		pkt := &response.SpawnNpc{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.RemoveNpc{}).Opcode():
		pkt := &response.RemoveNpc{}
		pkt.Deserialize(reader)
		return pkt, nil
	case (&response.Dialog{}).Opcode():
		if len(body) < 6 {
			return nil, fmt.Errorf("short dialog")
		}
		var pkt interface {
			Deserialize(reader *stream.StreamReader)
		}
		switch constant.DialogType(body[5]) {
		case constant.DialogTypeYesNo:
			pkt = &response.DialogYesNo{}
		case constant.DialogTypeInput:
			pkt = &response.DialogInput{}
		case constant.DialogTypeList:
			pkt = &response.DialogList{}
		case constant.DialogTypeStyle:
			pkt = &response.DialogStyle{}
		case constant.DialogTypeAccept, constant.DialogTypeAcceptEscape:
			pkt = &response.DialogAccept{}
		default:
			pkt = &response.Dialog{}
		}
		pkt.Deserialize(reader)
		return pkt, nil
	default:
		return nil, fmt.Errorf("unknown opcode %d", opcode)
	}
}
