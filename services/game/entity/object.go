package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type ObjectCore struct {
	self      Object
	OID       uint32
	Position  types.Vector2[int16]
	GameWorld GameWorld
	Map       *Map
	section   *section
	timers    map[string]*ObjectTimer
}

type ObjectBroadcastOption struct {
	SendRaw                  bool
	WithMe                   bool
	RecipientsRoleBelowPivot bool
}

func (obj *ObjectCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeObject
}

type Object interface {
	GetOID() uint32
	GetPK() uint32
	GetPosition() types.Vector2[int16]
	SetPosition(x, y int16)
	GetGameWorld() GameWorld
	GetMap() *Map
	getSection() *section
	setSection(section *section)
	GetObjectType() constant.ObjectType
	Is(typ constant.ObjectType) bool
	GetRole() constant.CharacterRole
	IsHidden() bool
	Send(p types.Packet, policy types.SendPolicy) error
	SendSpawnSyncToViewer(viewer *Character)
	SendDestroySyncToViewer(viewer *Character)
	Nears(filter constant.ObjectType, option *SearchOption) []Object
	ObjectsIn(filter constant.ObjectType, bounds types.Rect[int32]) []Object
	Broadcast(message types.Packet, option *ObjectBroadcastOption)
	BroadcastCall(fn func(Object), option *ObjectBroadcastOption)
	GetTimerEntry(key string) *ObjectTimer
	RemoveTimer(key string) bool
	RescheduleTimer(key string) bool
	SuspendTimers()
	ResumeTimers(pid *actor.PID)
}

func (obj *ObjectCore) GetOID() uint32 {
	return obj.OID
}

func (obj *ObjectCore) GetPK() uint32 {
	return obj.OID
}

func (obj *ObjectCore) GetPosition() types.Vector2[int16] {
	return obj.Position
}

func (obj *ObjectCore) SetPosition(x, y int16) {
	before := obj.Position
	obj.Position.X = x
	obj.Position.Y = y
	if obj.Map != nil && obj.self != nil {
		obj.Map.OnMoved(obj.self, before)
	}
}

func (obj *ObjectCore) GetGameWorld() GameWorld {
	return obj.GameWorld
}

func (obj *ObjectCore) Is(typ constant.ObjectType) bool {
	return obj.GetObjectType().Has(typ)
}

func (obj *ObjectCore) GetMap() *Map {
	return obj.Map
}

func (obj *ObjectCore) getSection() *section {
	return obj.section
}

func (obj *ObjectCore) setSection(section *section) {
	obj.section = section
}

func (obj *ObjectCore) GetRole() constant.CharacterRole {
	return constant.RoleUser
}

func (obj *ObjectCore) IsHidden() bool {
	return false
}

func (obj *ObjectCore) Send(types.Packet, types.SendPolicy) error {
	return nil
}

func (obj *ObjectCore) SendSpawnSyncToViewer(viewer *Character) {}

func (obj *ObjectCore) SendDestroySyncToViewer(viewer *Character) {}

func (o *ObjectCore) Nears(filter constant.ObjectType, option *SearchOption) []Object {
	pivot := o.self
	if pivot == nil {
		return nil
	}
	m := pivot.GetMap()
	if m == nil {
		return nil
	}
	out := make([]Object, 0)
	for _, cand := range m.GetObjectsNear(pivot.GetPosition(), filter, option) {
		if cand == nil || cand == pivot {
			continue
		}
		out = append(out, cand)
	}
	return out
}

func (o *ObjectCore) ObjectsIn(filter constant.ObjectType, bounds types.Rect[int32]) []Object {
	pivot := o.self
	if pivot == nil {
		return nil
	}
	m := pivot.GetMap()
	if m == nil || !bounds.Valid() {
		return nil
	}
	pos := pivot.GetPosition()
	world := bounds.AtOrigin(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)})
	out := make([]Object, 0)
	for _, obj := range m.GetObjectsIn(filter, world) {
		if obj == nil || obj == pivot {
			continue
		}
		out = append(out, obj)
	}
	return out
}

func (o *ObjectCore) BroadcastCall(fn func(Object), option *ObjectBroadcastOption) {
	if o.Map == nil || fn == nil {
		return
	}
	pivot := o.self
	if pivot == nil {
		return
	}
	withMe := false
	roleBelow := false
	if option != nil {
		withMe = option.WithMe
		roleBelow = option.RecipientsRoleBelowPivot
	}
	pivotRole := pivot.GetRole()
	for _, oc := range pivot.Nears(constant.ObjectTypeObject, nil) {
		if oc == nil {
			continue
		}
		if roleBelow && oc.GetRole() >= pivotRole {
			continue
		}
		fn(oc)
	}
	if withMe {
		fn(pivot)
	}
}

func (o *ObjectCore) Broadcast(message types.Packet, option *ObjectBroadcastOption) {
	if o.Map == nil {
		return
	}
	pivot := o.self
	if pivot == nil {
		return
	}
	policy := types.SEND_POLICY_ENCRYPT
	if option != nil && option.SendRaw {
		policy = types.SEND_POLICY_RAW
	}
	o.BroadcastCall(func(oc Object) {
		ch, ok := oc.(*Character)
		if !ok || ch == nil {
			return
		}
		_ = ch.Send(message, policy)
	}, option)
	if option != nil && option.WithMe {
		if _, ok := pivot.(*Character); !ok {
			_ = pivot.Send(message, policy)
		}
	}
}

func (fp *FieldPlacement) Is(typ constant.ObjectType) bool {
	return fp.GetObjectType().Has(typ)
}

var (
	_ Object = (*ObjectCore)(nil)
	_ Object = (*LifeCore)(nil)
	_ Object = (*Character)(nil)
	_ Object = (*Mob)(nil)
	_ Object = (*Npc)(nil)
	_ Object = (*Reactor)(nil)
	_ Object = (*FieldPlacement)(nil)
	_ Object = (*Mist)(nil)
	_ Object = (*Door)(nil)
	_ Object = (*Summon)(nil)
	_ Object = (*ItemCore)(nil)
	_ Object = (*Meso)(nil)
)
