package conn

import (
	"fmt"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
)

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
	default:
		return nil, fmt.Errorf("unknown opcode %d", opcode)
	}
}
