package response

import "github.com/boyism80/fm/stream"

// SnowballState shows both snowballs; the client reads the damage it reports only from the first one it gets.
type SnowballState struct {
	State     uint8
	SnowmanHP [2]int32
	Position  [2]int16
	Stage     [2]uint8
	Damage    [3]int16
}

func (p *SnowballState) Opcode() uint16 {
	return 0xD3
}

func (p *SnowballState) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(p.State)
	sw.Write32(p.SnowmanHP[0])
	sw.Write32(p.SnowmanHP[1])
	for i := range p.Position {
		sw.Write16(p.Position[i])
		sw.WriteU8(p.Stage[i])
	}
	for _, damage := range p.Damage {
		sw.Write16(damage)
	}
	return nil
}

func (p *SnowballState) Deserialize(reader *stream.StreamReader) {
	p.State = reader.ReadU8()
	p.SnowmanHP[0] = reader.Read32()
	p.SnowmanHP[1] = reader.Read32()
	for i := range p.Position {
		p.Position[i] = reader.Read16()
		p.Stage[i] = reader.ReadU8()
	}
	for i := range p.Damage {
		p.Damage[i] = reader.Read16()
	}
}

// SnowballHit shows the damage an attack did to a snowball or snowman.
type SnowballHit struct {
	Target uint8
	Damage int16
	Delay  int16
}

func (p *SnowballHit) Opcode() uint16 {
	return 0xD4
}

func (p *SnowballHit) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(p.Target)
	sw.Write16(p.Damage)
	sw.Write16(p.Delay)
	return nil
}

func (p *SnowballHit) Deserialize(reader *stream.StreamReader) {
	p.Target = reader.ReadU8()
	p.Damage = reader.Read16()
	p.Delay = reader.Read16()
}

// SnowballMessage shows a stage message for a team: 1-3 a stage passed, 4 its snowball stopped, 5 rolling again.
type SnowballMessage struct {
	Team    uint8
	Message uint8
}

func (p *SnowballMessage) Opcode() uint16 {
	return 0xD5
}

func (p *SnowballMessage) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(p.Team)
	sw.WriteU8(p.Message)
	return nil
}

func (p *SnowballMessage) Deserialize(reader *stream.StreamReader) {
	p.Team = reader.ReadU8()
	p.Message = reader.ReadU8()
}

// SnowballKnockback pushes the character back from a snowball it ran into.
type SnowballKnockback struct{}

func (p *SnowballKnockback) Opcode() uint16 {
	return 0xD6
}

func (p *SnowballKnockback) Serialize(sw *stream.StreamWriter) error {
	return nil
}

func (p *SnowballKnockback) Deserialize(reader *stream.StreamReader) {}
