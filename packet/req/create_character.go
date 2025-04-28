package req

import "github.com/boyism80/fm/stream"

type CreateCharacter struct {
	Name   string
	Face   uint32
	Hair   uint32
	Top    uint32
	Bottom uint32
	Shoes  uint32
	Weapon uint32
}

func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) error {
	name, _ := reader.ReadStr16()
	a.Name = name

	face, _ := reader.ReadU32()
	a.Face = face

	hair, _ := reader.ReadU32()
	a.Hair = hair

	top, _ := reader.ReadU32()
	a.Top = top

	bottom, _ := reader.ReadU32()
	a.Bottom = bottom

	shoes, _ := reader.ReadU32()
	a.Shoes = shoes

	weapon, _ := reader.ReadU32()
	a.Weapon = weapon

	return nil
}
