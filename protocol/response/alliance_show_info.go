package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceShowInfo struct {
	Info *dto.AllianceInfo
}

func (p *AllianceShowInfo) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceShowInfo) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CShowInfo))
	if p.Info == nil {
		w.WriteU8(0)
		return nil
	}
	w.WriteU8(1)
	writeAllianceInfo(w, p.Info)
	return nil
}

func (p *AllianceShowInfo) Deserialize(*stream.StreamReader) {}
