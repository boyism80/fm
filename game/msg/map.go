package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
	"github.com/boyism80/fm/game/protocol/resp"
)

// 맵에 오브젝트 등록
type EnterMap struct {
	ID         uint32
	PID        *actor.PID
	SpawnPoint uint8
	Init       bool
}

// 맵에서 오브젝트 제거
type LeaveMap struct {
	ID uint32
}

// 맵에 존재하는 모든 PID 요청
type QueryAllPIDs struct {
	ReplyTo *actor.PID
}

// 응답 메시지: PID 리스트
type MapPIDList struct {
	Targets []*actor.PID
}

type MapNotifyCharacterWarped struct {
	Sender *actor.PID
}

type MapBroadcast struct {
	Sender     *actor.PID
	Pivot      types.Vector2[int16]
	Message    any
	ExceptSelf bool
	Excepts    map[*actor.PID]struct{}
}

type MapSpawnItem struct {
	Item    entity.Item
	Owner   *actor.PID
	OwnerID uint32
}

type MapSpawnItems struct {
	Items    []entity.Dropable
	Position types.Vector2[int16]
	Owner    *actor.PID
	OwnerID  uint32
}

type MapSpawnMeso struct {
	Count        int32
	SpawnedPoint types.Vector2[int16]
	DestPoint    types.Vector2[int16]
	Owner        *actor.PID
	OwnerID      uint32
	DropType     constant.DropType
}

type MapItemLoot struct {
	Actor       *actor.PID
	OID         uint32
	CharacterId uint32
	Position    types.Vector2[int16]
}

type MapRemoveItem struct {
	Actor       *actor.PID
	OID         uint32
	CharacterId uint32
	Mode        resp.RemoveItemType
	Position    types.Vector2[int16]
}

type MapChange struct {
	CharacterId uint32
	Sender      *actor.PID
	To          *actor.PID
	SpawnPoint  uint8
}

type SendMessage struct {
	OID     uint32
	Message any
}

type MapRepeatSpawnMobs struct {
}

type MapDieMob struct {
	Sender        *actor.PID
	OID           uint32
	Position      types.Vector2[int16]
	AnimationType constant.MobDieAnimationType
}

type MapClearMobs struct {
	AnimationType constant.MobDieAnimationType
}

type MapMoveMob struct {
	req.MoveMob
	Sender *actor.PID
}

type MapSpawningMob struct {
	Position types.Vector2[int16]
	MobId    uint32
}

type MapSpawnedMob struct {
	PID *actor.PID
	OID uint32
}

type MapCharacterAttack struct {
	Sender      *actor.PID
	CharacterId uint32
	AttackInfo  protocol.AttackInfo
	Position    types.Vector2[int16]
}
