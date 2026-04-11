package entity

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Door struct {
	ObjectCore
	OwnerID          uint32
	SkillID          constant.SkillID
	OppositeMapID    uint32
	OppositePosition types.Vector2[int16]
	ExpiresAt        time.Time
}

func (d *Door) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeDoor
}

func (d *Door) Is(typ constant.ObjectType) bool {
	return d.GetObjectType().Has(typ)
}

func (d *Door) SendSpawnSyncToViewer(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	viewer.Send(&response.SpawnDoor{
		OwnerID:  d.OwnerID,
		Position: d.Position,
		Animated: false,
	}, types.SEND_POLICY_ENCRYPT)
	pos := d.Position
	viewer.Send(&response.SpawnPortal{
		TownMapID:   d.OppositeMapID,
		TargetMapID: uint32(d.Map.Wz.ID),
		SkillID:     uint32(d.SkillID),
		Position:    &pos,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) IsExpired(now time.Time) bool {
	return !d.ExpiresAt.IsZero() && !now.Before(d.ExpiresAt)
}

func (d *Door) Spawn(animated bool) {
	m := d.GetMap()
	if m == nil {
		return
	}
	m.AddDoor(d)
}

func (d *Door) Remove(animated bool) {
	m := d.GetMap()
	if m == nil || d.OID == 0 {
		return
	}
	m.RemoveDoor(d.OID, animated)
}
