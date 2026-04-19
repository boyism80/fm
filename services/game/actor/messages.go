package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type RequestSpawnDoor struct {
	ReplyTo        *actor.PID
	CharacterID    uint32
	OwnerID        uint32
	SkillID        constant.SkillID
	FieldMapID     uint32
	FieldPortalID  uint8
	PartyOwnerSlot int
	PartyID        *uint32
	FieldAnchor    types.Vector2[int16]
}

type ResponseSpawnDoor struct {
	Ok                 bool
	CharacterID        uint32
	OwnerID            uint32
	SkillID            constant.SkillID
	ReturnPortalID     uint8
	TownPortalPosition types.Vector2[int16]
	FieldMapID         uint32
	FieldPortalID      uint8
	PartyID            *uint32
	FieldAnchor        types.Vector2[int16]
}

type RemoveDoor struct {
	OwnerID uint32
	SkillID uint32
}

type ResumeLua struct {
	Root   *lua.LState
	Thread *lua.LState
}

type AddCharacter struct {
	Character  *entity.Character
	SpawnPoint uint8
	Init       bool
}

type RemoveCharacter struct {
	CharacterID uint32
}

type WarpCharacter struct {
	Character *entity.Character
	Portal    uint8
}

type TimerTick struct {
	HandlerName string
}

type SyncPartySnapshot struct {
	Snapshot *internal.PartySnapshot
}

type ClearPartyByPartyID struct {
	PartyID uint32
}

type SyncCharacterPartyState struct {
	CharacterID uint32
	PartyID     *uint32
}

type EnsureDeliver struct {
	CorrelationID    uint64
	CharacterID      uint32
	DeadlineUnixNano int64
	Caller           *actor.PID
	Inner            interface{}
}

type EnsureResult struct {
	CorrelationID uint64
	OK            bool
	Reason        string
}

type EnsureCoordinator interface {
	EnsureRedispatch(d *EnsureDeliver)
	EnsureComplete(correlationID uint64)
}

type DeliverPartyInvite struct {
	CharacterID uint32
	PartyID     uint32
	InviterName string
	PartySearch bool
}

type DeliverPartyStatusMessage struct {
	CharacterID uint32
	Code        pconst.PartyStatusCode
	Name        string
}

type DeliverPartyUpdateJoin struct {
	CharacterID uint32
	ForChannel  int32
	PartyID     uint32
	JoinName    string
	LeaderID    uint32
	Members     []response.PartyMemberStatus
}

type DeliverPartyUpdateLeave struct {
	CharacterID uint32
	ForChannel  int32
	PartyID     uint32
	TargetID    uint32
	TargetName  string
	LeaderID    uint32
	Members     []response.PartyMemberStatus
	Expelled    bool
}

type DeliverPartyUpdateDisband struct {
	CharacterID uint32
	PartyID     uint32
	LeaderID    uint32
}

type DeliverPartyUpdateLeaderChange struct {
	CharacterID          uint32
	NewLeaderCharacterID uint32
	ByDisconnect         bool
}

type DeliverPartyUpdateLogOnOff struct {
	CharacterID uint32
	ForChannel  int32
	PartyID     uint32
	LeaderID    uint32
	Members     []response.PartyMemberStatus
}

type DeliverPartyUpdateSilent struct {
	CharacterID uint32
	ForChannel  int32
	PartyID     uint32
	LeaderID    uint32
	Members     []response.PartyMemberStatus
}
