package request

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

func (*CreateCharacter) Opcode() byte { return 0x08 }

func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.Name)
	writer.WriteU32(a.Face)
	writer.WriteU32(a.Hair)
	writer.WriteU32(a.Top)
	writer.WriteU32(a.Bottom)
	writer.WriteU32(a.Shoes)
	writer.WriteU32(a.Weapon)
	return nil
}

func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) {
	a.Name = reader.ReadStr16()
	a.Face = reader.ReadU32()
	a.Hair = reader.ReadU32()
	a.Top = reader.ReadU32()
	a.Bottom = reader.ReadU32()
	a.Shoes = reader.ReadU32()
	a.Weapon = reader.ReadU32()

}
