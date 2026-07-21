package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type DamageMob struct {
	OID     uint32
	Display constant.MobDamageDisplayType
	Damage  int32
	HP      int32
	MaxHP   int32
}

func (p *DamageMob) Opcode() uint16 {
	return 0xB3
}

func (p *DamageMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.WriteU8(uint8(p.Display))
	writer.Write32(p.Damage)
	if p.Display.IncludesHpMaxHp() {
		writer.Write32(p.HP)
		writer.Write32(p.MaxHP)
	}
	return nil
}

func (p *DamageMob) Deserialize(reader *stream.StreamReader) {
}
