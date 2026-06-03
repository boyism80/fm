package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceDisband struct {
	AllianceID uint32
}

func (p *AllianceDisband) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceDisband) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CDisband))
	w.WriteU32(p.AllianceID)
	return nil
}

func (p *AllianceDisband) Deserialize(*stream.StreamReader) {}
