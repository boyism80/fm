package entity

import (
	"errors"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

var ErrMonsterBookCardNotOwned = errors.New("monster book card not owned")

func (ch *Character) RegisterCard(cardID uint32) {
	count := ch.MonsterBook.Cards[cardID]
	if count >= constant.MonsterBookCardMax {
		ch.Listener.OnMonsterBookCardFull(ch)
		return
	}

	count++
	ch.MonsterBook.Cards[cardID] = count
	ch.Listener.OnMonsterBookCardRegistered(ch, cardID, count)
	ch.Listener.OnShowSelfEffect(ch, response.EffectTypeRegisterCard)
}

func (ch *Character) SetMonsterBookCover(cardID uint32) error {
	if ch.MonsterBook.Cover == cardID {
		return nil
	}
	if cardID != 0 && ch.MonsterBook.Cards[cardID] == 0 {
		return ErrMonsterBookCardNotOwned
	}

	ch.MonsterBook.Cover = cardID
	ch.Listener.OnMonsterBookCover(ch, cardID)
	return nil
}

func (ch *Character) ResetMonsterBook() {
	ch.MonsterBook.Cards = make(map[uint32]uint32)
	ch.MonsterBook.Cover = 0
	ch.Listener.OnMonsterBookCover(ch, 0)
}
