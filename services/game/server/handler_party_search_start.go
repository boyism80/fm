package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type PartySearchStart struct {
	gs *GameServer
}

func (PartySearchStart) New(gs *GameServer) *PartySearchStart {
	return &PartySearchStart{
		gs: gs,
	}
}

func (h *PartySearchStart) Handle(ctx *core.ClientContext, req *request.PartySearchStart) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	if !h.canSearch(ch) {
		return nil
	}

	if req.MaxLevel <= req.MinLevel ||
		req.MaxLevel-req.MinLevel > 30 ||
		req.MembersNeeded <= 0 ||
		req.MembersNeeded > 6 ||
		req.MinLevel > int32(ch.GetLevel()) ||
		req.MaxLevel < int32(ch.GetLevel()) ||
		req.ClassMask == 0 {
		return nil
	}

	cfg := &game.PartySearchConfig{
		MinLevel:      req.MinLevel,
		MaxLevel:      req.MaxLevel,
		MembersNeeded: req.MembersNeeded,
		ClassMask:     req.ClassMask,
	}

	if pid := ch.Party().ID(); pid != nil {
		party := h.gs.systems.Get[*game.PartySystem]().Get(*pid)
		if party == nil || party.GetLeaderCharacterId() != ch.GetID() {
			return nil
		}
	}

	if ch.Party().ID() == nil {
		h.gs.systems.Get[*game.PartySystem]().CreateAsync(ctx.ActorContext, ch).Do(func(*internal.CreatePartyReply) error {
			if h.canSearch(ch) {
				ch.Party().SetSearch(cfg)
			}
			return nil
		}).OnError(func(err error) {
			log.Printf("PartySearchStart: %v", err)
		})
		return nil
	}
	ch.Party().SetSearch(cfg)
	return nil
}

func (h *PartySearchStart) canSearch(ch *game.Character) bool {
	if m := ch.GetMap(); m == nil || !ch.IsAlive() || m.BlocksTeleportStone() {
		ch.Listener().OnMessage(ch, constant.MsgPinkText, "이곳에서는 파티 찾기를 시도할 수 없습니다. 요청이 무시되었습니다.")
		return false
	}
	return true
}
