package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type DieMob struct {
	OID           uint32
	AnimationType constant.MobDieAnimationType
}

func (p *DieMob) Opcode() uint16 {
	return 0xAA
}

func (p *DieMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.WriteU8(uint8(p.AnimationType))
	switch p.AnimationType {
	case 4:
		writer.Write32(-1)
	}
	return nil
}

func (p *DieMob) Deserialize(reader *stream.StreamReader) {
}
