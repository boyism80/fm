package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceUpdateLeader struct {
	AllianceID  uint32
	OldLeaderID uint32
	NewLeaderID uint32
}

func (p *AllianceUpdateLeader) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceUpdateLeader) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CUpdateLeader))
	w.WriteU32(p.AllianceID)
	w.WriteU32(p.OldLeaderID)
	w.WriteU32(p.NewLeaderID)
	return nil
}

func (p *AllianceUpdateLeader) Deserialize(*stream.StreamReader) {}
