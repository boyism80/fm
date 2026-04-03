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

func (a *SummonAttack) Deserialize(reader *stream.StreamReader) error {
	var err error
	a.SummonOID, err = reader.ReadU32()
	if err != nil {
		return err
	}
	a.Tick, err = reader.ReadU32()
	if err != nil {
		return err
	}
	a.Animation, err = reader.ReadU8()
	if err != nil {
		return err
	}
	numAttacked, err := reader.ReadU8()
	if err != nil {
		return err
	}
	a.Damages = make([]dto.AttackPair, 0, numAttacked)
	for i := 0; i < int(numAttacked); i++ {
		oid, err := reader.ReadU32()
		if err != nil {
			return err
		}
		if err := reader.Skip(14); err != nil {
			return err
		}
		damage, err := reader.ReadU32()
		if err != nil {
			return err
		}
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
	return nil
}
