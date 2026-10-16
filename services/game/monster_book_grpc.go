package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewMonsterBookFromInternalProto(owner *Character, cover uint32, cards []*internal.MonsterBookCard) *MonsterBook {
	book := &MonsterBook{owner: owner, cover: cover, cards: make(map[uint32]uint32, len(cards))}
	for _, pb := range cards {
		model, ok := owner.itemModel(pb.GetCardId()).(*wz.Consume)
		if !ok || model.ConsumeType != constant.ConsumeTypeCard || pb.GetCount() == 0 {
			continue
		}
		book.cards[pb.GetCardId()] = min(pb.GetCount(), constant.MonsterBookCardMax)
	}
	return book
}

func (b *MonsterBook) BuildRPC() []*internal.MonsterBookCard {
	out := make([]*internal.MonsterBookCard, 0, len(b.cards))
	for cardID, count := range b.cards {
		out = append(out, &internal.MonsterBookCard{CardId: cardID, Count: count})
	}
	return out
}
