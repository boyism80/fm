package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MiniGameRecord struct {
	Type   uint32
	Wins   int32
	Ties   int32
	Losses int32
	Score  int32
}

func (r *MiniGameRecord) serialize(writer *stream.StreamWriter) {
	writer.WriteU32(r.Type)
	writer.Write32(r.Wins)
	writer.Write32(r.Ties)
	writer.Write32(r.Losses)
	writer.Write32(r.Score)
}

func (r *MiniGameRecord) deserialize(reader *stream.StreamReader) {
	r.Type = reader.ReadU32()
	r.Wins = reader.Read32()
	r.Ties = reader.Read32()
	r.Losses = reader.Read32()
	r.Score = reader.Read32()
}

type MiniGameMember struct {
	MiniRoomVisitor
	Record MiniGameRecord
}

type MiniGameEnter struct {
	Type    uint8
	MySlot  uint8
	Members []MiniGameMember
	Title   string
	Piece   uint8
}

func (p *MiniGameEnter) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameEnter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultEnter))
	writer.WriteU8(p.Type)
	writer.WriteU8(pconst.MiniRoomGameUsers)
	writer.WriteU8(p.MySlot)
	for _, member := range p.Members {
		writer.WriteU8(member.Slot)
		member.Character.SerializeLook(writer)
		writer.WriteStr16(member.Character.Name)
	}
	writer.WriteU8(0xFF)
	for _, member := range p.Members {
		writer.WriteU8(member.Slot)
		member.Record.serialize(writer)
	}
	writer.WriteU8(0xFF)
	writer.WriteStr16(p.Title)
	writer.WriteU16(uint16(p.Piece))
	return nil
}

func (p *MiniGameEnter) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Type = reader.ReadU8()
	reader.ReadU8()
	p.MySlot = reader.ReadU8()
	for slot := int8(reader.ReadU8()); slot >= 0; slot = int8(reader.ReadU8()) {
		member := MiniGameMember{MiniRoomVisitor: MiniRoomVisitor{Slot: uint8(slot), Character: &dto.Character{}}}
		member.Character.DeserializeLook(reader)
		member.Character.Name = reader.ReadStr16()
		p.Members = append(p.Members, member)
	}
	for slot := int8(reader.ReadU8()); slot >= 0; slot = int8(reader.ReadU8()) {
		for i := range p.Members {
			if p.Members[i].Slot == uint8(slot) {
				p.Members[i].Record.deserialize(reader)
			}
		}
	}
	p.Title = reader.ReadStr16()
	p.Piece = uint8(reader.ReadU16())
}

type MiniGameVisited struct {
	MiniGameMember
}

func (p *MiniGameVisited) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameVisited) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultVisitor))
	writer.WriteU8(p.Slot)
	p.Character.SerializeLook(writer)
	writer.WriteStr16(p.Character.Name)
	p.Record.serialize(writer)
	return nil
}

func (p *MiniGameVisited) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Slot = reader.ReadU8()
	p.Character = &dto.Character{}
	p.Character.DeserializeLook(reader)
	p.Character.Name = reader.ReadStr16()
	p.Record.deserialize(reader)
}

type MiniGameReady struct {
	Ready bool
}

func (p *MiniGameReady) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameReady) Serialize(writer *stream.StreamWriter) error {
	result := pconst.MiniRoomResultUnready
	if p.Ready {
		result = pconst.MiniRoomResultReady
	}
	writer.WriteU8(uint8(result))
	return nil
}

func (p *MiniGameReady) Deserialize(reader *stream.StreamReader) {
	p.Ready = pconst.MiniRoomResult(reader.ReadU8()) == pconst.MiniRoomResultReady
}

type MiniGameExitAfter struct {
	Reserved bool
}

func (p *MiniGameExitAfter) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameExitAfter) Serialize(writer *stream.StreamWriter) error {
	result := pconst.MiniRoomResultCancelExit
	if p.Reserved {
		result = pconst.MiniRoomResultExitAfterGame
	}
	writer.WriteU8(uint8(result))
	return nil
}

func (p *MiniGameExitAfter) Deserialize(reader *stream.StreamReader) {
	p.Reserved = pconst.MiniRoomResult(reader.ReadU8()) == pconst.MiniRoomResultExitAfterGame
}

type MiniGameStart struct {
	Loser uint8
	Cards []uint32
}

func (p *MiniGameStart) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameStart) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultStart))
	writer.WriteU8(p.Loser)
	if p.Cards == nil {
		return nil
	}
	writer.WriteU8(uint8(len(p.Cards)))
	for _, card := range p.Cards {
		writer.WriteU32(card)
	}
	return nil
}

func (p *MiniGameStart) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Loser = reader.ReadU8()
	if reader.Remaining() == 0 {
		return
	}
	count := int(reader.ReadU8())
	p.Cards = make([]uint32, 0, count)
	for i := 0; i < count; i++ {
		p.Cards = append(p.Cards, reader.ReadU32())
	}
}

type MiniGameSkip struct {
	Slot uint8
}

func (p *MiniGameSkip) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameSkip) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultSkip))
	writer.WriteU8(p.Slot)
	return nil
}

func (p *MiniGameSkip) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Slot = reader.ReadU8()
}

type MiniGameRequestTie struct{}

func (p *MiniGameRequestTie) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameRequestTie) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultRequestTie))
	return nil
}

func (p *MiniGameRequestTie) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
}

type MiniGameDenyTie struct{}

func (p *MiniGameDenyTie) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameDenyTie) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultDenyTie))
	return nil
}

func (p *MiniGameDenyTie) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
}

type MiniGameMoveOmok struct {
	X     int32
	Y     int32
	Stone uint8
}

func (p *MiniGameMoveOmok) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameMoveOmok) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultMoveOmok))
	writer.Write32(p.X)
	writer.Write32(p.Y)
	writer.WriteU8(p.Stone)
	return nil
}

func (p *MiniGameMoveOmok) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.X = reader.Read32()
	p.Y = reader.Read32()
	p.Stone = reader.ReadU8()
}

type MiniGameSelectCard struct {
	FirstPick bool
	Card      uint8
	FirstCard uint8
	Result    uint8
}

func (p *MiniGameSelectCard) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameSelectCard) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultSelectCard))
	writer.WriteBoolean(p.FirstPick)
	writer.WriteU8(p.Card)
	if p.FirstPick {
		return nil
	}
	writer.WriteU8(p.FirstCard)
	writer.WriteU8(p.Result)
	return nil
}

func (p *MiniGameSelectCard) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.FirstPick = reader.ReadBool()
	p.Card = reader.ReadU8()
	if p.FirstPick {
		return
	}
	p.FirstCard = reader.ReadU8()
	p.Result = reader.ReadU8()
}

type MiniGameOver struct {
	Outcome pconst.MiniGameOutcome
	Winner  uint8
	Records []MiniGameRecord
}

func (p *MiniGameOver) Opcode() uint16 {
	return 0xEF
}

func (p *MiniGameOver) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultGameOver))
	writer.WriteU8(uint8(p.Outcome))
	if p.Outcome != pconst.MiniGameTie {
		writer.WriteU8(p.Winner)
	}
	for _, record := range p.Records {
		record.serialize(writer)
	}
	return nil
}

func (p *MiniGameOver) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Outcome = pconst.MiniGameOutcome(reader.ReadU8())
	if p.Outcome != pconst.MiniGameTie {
		p.Winner = reader.ReadU8()
	}
	for reader.Remaining() > 0 {
		var record MiniGameRecord
		record.deserialize(reader)
		p.Records = append(p.Records, record)
	}
}
