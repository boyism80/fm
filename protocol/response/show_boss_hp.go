package response

import "github.com/boyism80/fm/stream"

type ShowBossHp struct {
	MobID      uint32
	CurrentHP  int32
	MaxHP      int32
	TagColor   uint8
	TagBgColor uint8
}

func (p *ShowBossHp) Opcode() uint16 {
	return 0x5F
}

func (p *ShowBossHp) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(5)
	sw.WriteU32(p.MobID)
	sw.Write32(p.CurrentHP)
	sw.Write32(p.MaxHP)
	sw.WriteU8(p.TagColor)
	sw.WriteU8(p.TagBgColor)
	return nil
}

func (p *ShowBossHp) Deserialize(*stream.StreamReader) {
}
