package entity

import "github.com/boyism80/fm/services/game/constant"

type MonsterBook struct {
	Cover uint32
	Cards map[uint32]uint32
}

func (b *MonsterBook) Count() (normal uint32, special uint32) {
	for cardID := range b.Cards {
		if constant.IsSpecialMonsterCard(cardID) {
			special++
		} else {
			normal++
		}
	}
	return normal, special
}

func (b *MonsterBook) Level() uint32 {
	total := uint32(len(b.Cards))
	for i, limit := range []uint32{10, 30, 60, 100, 150, 210, 280} {
		if total <= limit {
			return uint32(i + 1)
		}
	}
	return 8
}
