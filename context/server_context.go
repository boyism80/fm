package context

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/data"
)

// 읽기 전용 인터페이스
type IServerContext interface {
	GameData() *data.GameData
	MapActor(mapId uint32) *protoactor.PID
}

// 실제 구조체
type ServerContext struct {
	gameData  *data.GameData
	mapActors map[uint32]*protoactor.PID
}

// 생성자
func NewServerContext(
	gameData *data.GameData,
	mapActors map[uint32]*protoactor.PID,
	serverCtxActor *protoactor.PID) *ServerContext {

	return &ServerContext{
		gameData:  gameData,
		mapActors: mapActors,
	}
}

// 인터페이스 구현
func (c *ServerContext) GameData() *data.GameData {
	return c.gameData
}

func (c *ServerContext) MapActor(mapId uint32) *protoactor.PID {
	return c.mapActors[mapId]
}

func (c *ServerContext) Load(id uint32, pid *protoactor.PID) {
	c.mapActors[id] = pid
}
