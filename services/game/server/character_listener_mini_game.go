package server

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/domain"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnMiniGameEntered(ch *domain.Character, g *domain.MiniGame) {
	slot, _ := g.SlotOf(ch)
	enter := &response.MiniGameEnter{
		Type:   g.Kind,
		MySlot: slot,
		Title:  g.Title,
		Piece:  g.Piece,
	}
	for _, member := range g.Members() {
		memberSlot, _ := g.SlotOf(member)
		enter.Members = append(enter.Members, response.MiniGameMember{
			MiniRoomVisitor: response.MiniRoomVisitor{Slot: memberSlot, Character: member.BuildS2C()},
			Record:          g.Record(member),
		})
	}
	ch.Send(enter, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameVisited(ch *domain.Character, g *domain.MiniGame, visitor *domain.Character) {
	ch.Send(&response.MiniGameVisited{MiniGameMember: response.MiniGameMember{
		MiniRoomVisitor: response.MiniRoomVisitor{Slot: 1, Character: visitor.BuildS2C()},
		Record:          g.Record(visitor),
	}}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameReady(ch *domain.Character, ready bool) {
	ch.Send(&response.MiniGameReady{Ready: ready}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameExitAfter(ch *domain.Character, reserved bool) {
	ch.Send(&response.MiniGameExitAfter{Reserved: reserved}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameStarted(ch *domain.Character, g *domain.MiniGame) {
	ch.Send(&response.MiniGameStart{Loser: 1 - g.Turn, Cards: g.Cards}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameSkipped(ch *domain.Character, slot uint8) {
	ch.Send(&response.MiniGameSkip{Slot: slot}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameTieRequested(ch *domain.Character) {
	ch.Send(&response.MiniGameRequestTie{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameTieDenied(ch *domain.Character) {
	ch.Send(&response.MiniGameDenyTie{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameOmokMoved(ch *domain.Character, x int32, y int32, stone uint8) {
	ch.Send(&response.MiniGameMoveOmok{X: x, Y: y, Stone: stone}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameCardSelected(ch *domain.Character, firstPick bool, card uint8, firstCard uint8, result uint8) {
	ch.Send(&response.MiniGameSelectCard{
		FirstPick: firstPick,
		Card:      card,
		FirstCard: firstCard,
		Result:    result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniGameOver(ch *domain.Character, g *domain.MiniGame, outcome pconst.MiniGameOutcome, winner uint8) {
	over := &response.MiniGameOver{Outcome: outcome, Winner: winner}
	for _, member := range g.Members() {
		over.Records = append(over.Records, g.Record(member))
	}
	ch.Send(over, types.SEND_POLICY_ENCRYPT)
}
