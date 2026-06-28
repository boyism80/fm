package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type Reactor struct {
	ObjectCore
	Wz                 *wz.Reactor
	Spawn              *ReactorSpawn
	State              byte
	TimerActive        bool
	TriggerCharacterID uint32
}

func (r *Reactor) GetTrigger() *Character {
	if r == nil || r.Map == nil || r.TriggerCharacterID == 0 {
		return nil
	}
	return r.Map.GetPlayer(r.TriggerCharacterID)
}

func (r *Reactor) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeReactor
}

func (r *Reactor) Is(typ constant.ObjectType) bool {
	return r.GetObjectType().Has(typ)
}

func (r *Reactor) ToDTO() *dto.Reactor {
	if r == nil {
		return nil
	}

	reactorID := uint32(0)
	facing := uint8(0)
	name := ""
	if r.Wz != nil {
		reactorID = r.Wz.ID
	}
	if r.Spawn != nil && r.Spawn.Wz != nil {
		facing = uint8(r.Spawn.Wz.FacingDirection)
		name = r.Spawn.Wz.Name
	}

	return &dto.Reactor{
		OID:       r.OID,
		ReactorID: reactorID,
		State:     r.State,
		Position:  r.Position,
		Facing:    facing,
		Name:      name,
	}
}

func (r *Reactor) SendSpawnSyncToViewer(viewer *Character) {
	if r == nil || viewer == nil {
		return
	}

	viewer.Send(&response.SpawnReactor{
		Reactor: r.ToDTO(),
	}, types.SEND_POLICY_ENCRYPT)
}
