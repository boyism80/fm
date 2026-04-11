package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SummonAttack struct {
	SummonOID uint32
	Tick      uint32
	Animation uint8
	Damages   []dto.AttackPair
}

func (a *SummonAttack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *SummonAttack) Deserialize(reader *stream.StreamReader) {
	a.SummonOID = reader.ReadU32()
	a.Tick = reader.ReadU32()
	a.Animation = reader.ReadU8()
	numAttacked := reader.ReadU8()
	a.Damages = make([]dto.AttackPair, 0, numAttacked)
	for i := 0; i < int(numAttacked); i++ {
		oid := reader.ReadU32()
		reader.Skip(14)
		damage := reader.ReadU32()
		a.Damages = append(a.Damages, dto.AttackPair{
			OID: oid,
			DamagePairs: []dto.DamagePair{
				{
					Damage:  damage,
					Unknown: false,
				},
			},
		})
	}
}
