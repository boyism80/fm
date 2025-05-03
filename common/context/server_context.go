package context

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/data"
)

// 실제 구조체
type ServerContext struct {
	GameData  *data.GameData
	MapActors map[uint32]*protoactor.PID
}

// 생성자
func NewServerContext(
	GameData *data.GameData,
	mapActors map[uint32]*protoactor.PID,
	serverCtxActor *protoactor.PID) *ServerContext {

	return &ServerContext{
		GameData:  GameData,
		MapActors: mapActors,
	}
}
