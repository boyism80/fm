package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MoveFragment interface {
	GetStance() uint8
	Serialize(sw *stream.StreamWriter)
	Deserialize(reader *stream.StreamReader)
}

type BasicMovement struct {
	Command uint8
	Stance  uint8
}

type AbsoluteLifeMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Duration int16
	Velocity types.Vector2[int16]
	Offset   types.Vector2[int16]
	Foothold int16
}

type AranMovement struct {
	*BasicMovement
	Duration int16
}

type ChairMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Foothold int16
	Duration int16
}

type JumpDownMovement struct {
	*BasicMovement
	Position     types.Vector2[int16]
	Duration     int16
	Velocity     types.Vector2[int16]
	Foothold     int16
	FallFoothold int16
}

type NoneMovement struct {
	*BasicMovement
}

type RelativeLifeMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Duration int16
}

type TeleportMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Foothold int16
	Duration int16
}

func (m *BasicMovement) GetStance() uint8 {
	return m.Stance
}

func (m *AbsoluteLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.Write16(m.Velocity.Y)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *AbsoluteLifeMovement) Deserialize(reader *stream.StreamReader) {
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Velocity.X = reader.Read16()
	m.Velocity.Y = reader.Read16()
	m.Foothold = reader.Read16()
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func (m *AranMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *AranMovement) Deserialize(reader *stream.StreamReader) {
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func (m *ChairMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *ChairMovement) Deserialize(reader *stream.StreamReader) {
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Foothold = reader.Read16()
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func (m *JumpDownMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.Write16(m.Velocity.Y)
	writer.Write16(m.Foothold)
	writer.Write16(m.FallFoothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *JumpDownMovement) Deserialize(reader *stream.StreamReader) {
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Velocity.X = reader.Read16()
	m.Velocity.Y = reader.Read16()
	m.Foothold = reader.Read16()
	m.FallFoothold = reader.Read16()
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func (m *NoneMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.WriteU8(m.Stance)
}

func (m *NoneMovement) Deserialize(reader *stream.StreamReader) {
	m.Stance = reader.ReadU8()
}

func (m *RelativeLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *RelativeLifeMovement) Deserialize(reader *stream.StreamReader) {
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func (m *TeleportMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *TeleportMovement) Deserialize(reader *stream.StreamReader) {
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Foothold = reader.Read16()
	m.Stance = reader.ReadU8()
	m.Duration = reader.Read16()
}

func ReadMovements(reader *stream.StreamReader) []MoveFragment {
	numCommands := reader.ReadU8()

	fragments := make([]MoveFragment, 0, numCommands)

	for range int(numCommands) {
		basic := &BasicMovement{Command: reader.ReadU8()}

		var frag MoveFragment
		switch basic.Command {
		case 0, 5, 17:
			frag = &AbsoluteLifeMovement{BasicMovement: basic}
		case 15:
			frag = &JumpDownMovement{BasicMovement: basic}
		case 1, 2, 6, 12, 13, 16, 18:
			frag = &RelativeLifeMovement{BasicMovement: basic}
		case 3, 4, 7, 8, 9, 11:
			frag = &TeleportMovement{BasicMovement: basic}
		case 10:
			frag = &NoneMovement{BasicMovement: basic}
		case 14:
			frag = &ChairMovement{BasicMovement: basic}
		default:
			frag = &AranMovement{BasicMovement: basic}
		}
		frag.Deserialize(reader)
		fragments = append(fragments, frag)
	}

	return fragments
}
