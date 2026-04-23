package request

import "github.com/boyism80/fm/stream"

type DamageSummon struct {
	OID           uint32
	Unknown       uint8
	Damage        uint32
	MonsterIdFrom uint32
}

func (*DamageSummon) Opcode() byte { return 0x8E }

func (d *DamageSummon) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (d *DamageSummon) Deserialize(reader *stream.StreamReader) {
	d.OID = reader.ReadU32()
	d.Unknown = reader.ReadU8()
	d.Damage = reader.ReadU32()
	d.MonsterIdFrom = reader.ReadU32()
}
