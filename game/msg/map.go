package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
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

type MapSpawnNpc struct {
	Sender *actor.PID
}

type MapBroadcastRange struct {
	Sender     *actor.PID
	Pivot      types.Vector2[int16]
	Message    any
	ExceptSelf bool
}

type MapSpawnItem struct {
	Item    entity.Item
	Owner   *actor.PID
	OwnerID uint32
}

type MapSpawnMeso struct {
	Count        int32
	SpawnedPoint types.Vector2[int16]
	Owner        *actor.PID
	OwnerID      uint32
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
