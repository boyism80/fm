package dao

import "github.com/boyism80/fm/stream"

type MonsterBook struct {
	Cards map[uint32]uint32 // 카드ID -> 레벨 or 상태 (ex: 획득 여부 등)
}

func (mb *MonsterBook) Serialize(writer *stream.StreamWriter) {
	writer.WriteU16(uint16(len(mb.Cards)))

	for cardId := range mb.Cards {
		writer.WriteU16(uint16(cardId))
	}
}
