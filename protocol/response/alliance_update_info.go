package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceUpdateInfo struct {
	Info *dto.AllianceInfo
}

func (p *AllianceUpdateInfo) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceUpdateInfo) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CUpdateInfo))
	if p.Info != nil {
		writeAllianceInfo(w, p.Info)
	}
	return nil
}

func (p *AllianceUpdateInfo) Deserialize(*stream.StreamReader) {}
