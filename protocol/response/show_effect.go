package response

import (
	"github.com/boyism80/fm/stream"
)

type EffectType uint8

const (
	EffectTypeLevelUp         EffectType = 0
	EffectTypeJobChange       EffectType = 8
	EffectTypeQuestCompletion EffectType = 9
	EffectTypeRegisterCard    EffectType = 13
	EffectTypeItemLevelUp     EffectType = 15
)

type ShowEffect struct {
	CharacterID uint32
	Type        EffectType
}

func (p *ShowEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(uint8(p.Type))
	return nil
}

func (p *ShowEffect) Deserialize(reader *stream.StreamReader) {
}
