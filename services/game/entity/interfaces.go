package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type GameWorld interface {
	core.Server
	GetResources() *wz.Resources
	GetMap(mapId uint32) *Map
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	SaveCharactersAsync(ctx actor.Context, chars []*Character) *async.Promise
	RequestWarp(character *Character, targetMap *Map, spawnPoint uint8) error
	DispatchRunCharacterTimer(pid *actor.PID, payload *c_actor.RunCharacterTimer)
	NotifyDoorRemove(ownerID uint32, skillID uint32, counterpartMapWZID uint32)
	RequestSpawnReturnMapDoor(ch *Character, skillID constant.SkillID)
	GetPartyByID(partyID uint32) *Party
	RequestAutoInvitePartyAsync(ctx actor.Context, inviterCharacterID uint32, targetCharacterIDs []uint32) *async.Promise
}
