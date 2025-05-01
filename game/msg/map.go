package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
)

// 맵에 오브젝트 등록
type EnterMap struct {
	Id  uint32
	PID *protoactor.PID
}

// 맵에서 오브젝트 제거
type LeaveMap struct {
	Id uint32
}

// 맵에 존재하는 모든 PID 요청
type QueryAllPIDs struct {
	ReplyTo *protoactor.PID
}

// 응답 메시지: PID 리스트
type MapPidList struct {
	Targets []*protoactor.PID
}

type MapBroadcastRange struct {
	Sender  *actor.PID
	Pivot   types.Vec2
	Message interface{}
}
