package game

import (
	"errors"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
)

var ErrMonsterBookCardNotOwned = errors.New("monster book card not owned")

type MonsterBook struct {
	owner *Character
	cover uint32
	cards map[uint32]uint32
}

func (b *MonsterBook) Count() (normal uint32, special uint32) {
	for cardID, count := range b.cards {
		model := b.owner.itemModel(cardID)
		if model != nil && model.GetTraits()&constant.ItemTraitSpecialMonsterCard != 0 {
			special += count
		} else {
			normal += count
		}
	}
	return normal, special
}

func (b *MonsterBook) Level() uint32 {
	normal, special := b.Count()
	total := normal + special
	for i, limit := range []uint32{0, 10, 30, 60, 100, 150, 210, 280} {
		if total <= limit {
			return uint32(i + 1)
		}
	}
	return 8
}

func (b *MonsterBook) Register(cardID uint32) {
	count := b.cards[cardID]
	if count >= constant.MonsterBookCardMax {
		b.owner.listener.OnMonsterBookCardFull(b.owner)
		return
	}

	count++
	b.cards[cardID] = count
	b.owner.listener.OnMonsterBookCardRegistered(b.owner, cardID, count)
	b.owner.listener.OnShowSelfEffect(b.owner, response.EffectTypeRegisterCard)
}

func (b *MonsterBook) SetCover(cardID uint32) error {
	if b.cover == cardID {
		return nil
	}
	if cardID != 0 && b.cards[cardID] == 0 {
		return ErrMonsterBookCardNotOwned
	}

	b.cover = cardID
	b.owner.listener.OnMonsterBookCover(b.owner, cardID)
	return nil
}

func (b *MonsterBook) Reset() {
	b.cards = make(map[uint32]uint32)
	b.cover = 0
	b.owner.listener.OnMonsterBookCover(b.owner, 0)
}
