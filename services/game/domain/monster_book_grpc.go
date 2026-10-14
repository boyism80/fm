package domain

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func NewMonsterBookFromInternalProto(owner *Character, cover uint32, cards []*internal.MonsterBookCard) *MonsterBook {
	book := &MonsterBook{owner: owner, Cover: cover, Cards: make(map[uint32]uint32, len(cards))}
	for _, pb := range cards {
		model, ok := owner.itemModel(pb.GetCardId()).(*wz.Consume)
		if !ok || model.ConsumeType != constant.ConsumeTypeCard || pb.GetCount() == 0 {
			continue
		}
		book.Cards[pb.GetCardId()] = min(pb.GetCount(), constant.MonsterBookCardMax)
	}
	return book
}

func (b *MonsterBook) BuildRPC() []*internal.MonsterBookCard {
	out := make([]*internal.MonsterBookCard, 0, len(b.Cards))
	for cardID, count := range b.Cards {
		out = append(out, &internal.MonsterBookCard{CardId: cardID, Count: count})
	}
	return out
}
