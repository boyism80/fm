package msg

import "github.com/boyism80/fm/common/types"

type InvokeHandler struct {
	Opcode uint16
	Data   []byte
}

type SendProtocol struct {
	Protocol types.Packet
	Policy   types.SendPolicy
}
