package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

func reportPartyMemberSnapshotRequest(worldID uint32, ch *entity.Character) (*internal.ReportPartyMemberSnapshotRequest, bool) {
	if ch == nil {
		return nil, false
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return nil, false
	}
	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.GetMapID()
	}
	req := &internal.ReportPartyMemberSnapshotRequest{
		WorldId:     worldID,
		CharacterId: ch.GetID(),
		Level:       uint32(ch.GetLevel()),
		ClassId:     uint32(ch.Class),
		MapId:       mapID,
	}
	if ds := ch.GetDoors(); len(ds) > 0 && ds[0] != nil {
		d := ds[0]
		req.Door = &internal.PartyDoor{
			Town:   d.ReturnMapID,
			Target: d.FieldMapID,
			X:      int32(d.Position.X),
			Y:      int32(d.Position.Y),
		}
	}
	return req, true
}

func (gs *GameServer) ReportPartyMemberSnapshotAsync(ch *entity.Character) {
	if gs == nil || ch == nil || gs.internalClient == nil {
		return
	}
	req, ok := reportPartyMemberSnapshotRequest(uint32(gs.config.WorldId), ch)
	if !ok {
		return
	}
	cid := ch.GetID()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		defer cancel()
		if _, err := gs.internalClient.ReportPartyMemberSnapshot(ctx, req); err != nil {
			log.Printf("ReportPartyMemberSnapshot async char %d: %v", cid, err)
		}
	}()
}
