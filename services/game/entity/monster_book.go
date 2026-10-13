package entity

import (
	"errors"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

var ErrMonsterBookCardNotOwned = errors.New("monster book card not owned")

type MonsterBook struct {
	owner *Character
	Cover uint32
	Cards map[uint32]uint32
}

func (b *MonsterBook) Count() (normal uint32, special uint32) {
	for cardID := range b.Cards {
		model := b.owner.itemModel(cardID)
		if model != nil && model.GetTraits()&constant.ItemTraitSpecialMonsterCard != 0 {
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

func (b *MonsterBook) Register(cardID uint32) {
	count := b.Cards[cardID]
	if count >= constant.MonsterBookCardMax {
		b.owner.Listener.OnMonsterBookCardFull(b.owner)
		return
	}

	count++
	b.Cards[cardID] = count
	b.owner.Listener.OnMonsterBookCardRegistered(b.owner, cardID, count)
	b.owner.Listener.OnShowSelfEffect(b.owner, response.EffectTypeRegisterCard)
}

func (b *MonsterBook) SetCover(cardID uint32) error {
	if b.Cover == cardID {
		return nil
	}
	if cardID != 0 && b.Cards[cardID] == 0 {
		return ErrMonsterBookCardNotOwned
	}

	b.Cover = cardID
	b.owner.Listener.OnMonsterBookCover(b.owner, cardID)
	return nil
}

func (b *MonsterBook) Reset() {
	b.Cards = make(map[uint32]uint32)
	b.Cover = 0
	b.owner.Listener.OnMonsterBookCover(b.owner, 0)
}
