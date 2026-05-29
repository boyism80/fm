package response

import "github.com/boyism80/fm/stream"

type DamageMob struct {
	OID    uint32
	Damage int32
}

func (p *DamageMob) Opcode() uint16 {
	return 0xAC
}

func (p *DamageMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Damage)
	return nil
}

func (p *DamageMob) Deserialize(reader *stream.StreamReader) {
}
