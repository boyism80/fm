package request

import "github.com/boyism80/fm/stream"

type DamageSummon struct {
	OID           uint32
	Unknown       uint8
	Damage        uint32
	MonsterIdFrom uint32
}

func (d *DamageSummon) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (d *DamageSummon) Deserialize(reader *stream.StreamReader) error {
	var err error
	d.OID, err = reader.ReadU32()
	if err != nil {
		return err
	}
	d.Unknown, err = reader.ReadU8()
	if err != nil {
		return err
	}
	d.Damage, err = reader.ReadU32()
	if err != nil {
		return err
	}
	d.MonsterIdFrom, err = reader.ReadU32()
	if err != nil {
		return err
	}
	return nil
}
