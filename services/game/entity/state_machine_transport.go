package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/response"
)

func (sm *StateMachine) WarpAll(ctx actor.Context, fromMapID, toMapID uint32, spawnPoint uint8) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil {
		return
	}
	ms := sm.Group.GameWorld.GetMapSystem()
	if ms == nil {
		return
	}
	from := sm.Map(fromMapID)
	if from == nil {
		return
	}
	to := ms.Find(sm, toMapID)
	if to == nil {
		return
	}
	players := make([]*Character, 0)
	for _, obj := range from.GetAllPlayers() {
		ch, ok := obj.(*Character)
		if !ok || ch == nil || sm.HasPlayer(ch) == false {
			continue
		}
		players = append(players, ch)
	}
	for _, ch := range players {
		if ch.GetHp() == 0 {
			returnID := uint32(0)
			if from.Wz != nil {
				returnID = uint32(from.Wz.ReturnMapId)
			}
			if returnID == 0 || returnID == fromMapID {
				continue
			}
			ret := ms.Find(sm, returnID)
			if ret == nil {
				continue
			}
			_ = ch.Warp(ctx, ret, 0)
			continue
		}
		_ = ch.Warp(ctx, to, spawnPoint)
	}
}

func (sm *StateMachine) BroadcastShip(mapID uint32, effect uint16) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil {
		return
	}
	ms := sm.Group.GameWorld.GetMapSystem()
	m := ms.Find(sm, mapID)
	if m == nil {
		return
	}
	if effect == response.ShipSpecialBalrog {
		ms.Broadcast(m, &response.ShipSpecialEffect{Effect: effect})
		return
	}
	ms.Broadcast(m, &response.ShipState{State: effect})
}
