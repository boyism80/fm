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
	Wz          *wz.Reactor
	Spawn       *ReactorSpawn
	State       byte
	TimerActive bool
}

func (r *Reactor) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeReactor
}

func (r *Reactor) Is(typ constant.ObjectType) bool {
	return r.GetObjectType().Has(typ)
}

func (r *Reactor) ReactorID() uint32 {
	if r == nil || r.Spawn == nil || r.Spawn.Wz == nil {
		return 0
	}
	return r.Spawn.Wz.ReactorID
}

func (r *Reactor) ToDTO() *dto.Reactor {
	if r == nil {
		return nil
	}

	reactorID := uint32(0)
	facing := uint8(0)
	name := ""
	if r.Spawn != nil && r.Spawn.Wz != nil {
		reactorID = r.Spawn.Wz.ReactorID
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
