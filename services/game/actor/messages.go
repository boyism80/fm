package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type RequestSpawnDoor struct {
	ReplyTo        *actor.PID
	TargetMap      *entity.Map
	Key            entity.DoorKey
	CharacterID    uint32
	OwnerID        uint32
	SkillID        constant.SkillID
	Field          entity.DoorEndpoint
	PartyOwnerSlot int
	PartyID        *uint32
}

type ResponseSpawnDoor struct {
	Ok          bool
	Key         entity.DoorKey
	CharacterID uint32
	OwnerID     uint32
	SkillID     constant.SkillID
	Return      entity.DoorEndpoint
	Field       entity.DoorEndpoint
	PartyID     *uint32
}

type DespawnDoor struct {
	Map               *entity.Map
	Key               entity.DoorKey
	Animated          bool
	NotifyCounterpart bool
}

type SetDoorPartyID struct {
	Map     *entity.Map
	Key     entity.DoorKey
	PartyID *uint32
}

type ResumeLua struct {
	Root   *lua.LState
	Thread *lua.LState
	Args   []lua.LValue
}

type MapCallFunc func(ctx actor.Context, a *GameLogicActor) []lua.LValue

type MapCallAsyncFunc func(ctx actor.Context, a *GameLogicActor) *async.Promise[[]lua.LValue]

type MapCall struct {
	Run     MapCallFunc
	ReplyTo *actor.PID
	Root    *lua.LState
	Thread  *lua.LState
	Values  []lua.LValue
}

type MapCallAsync struct {
	Run     MapCallAsyncFunc
	ReplyTo *actor.PID
	Root    *lua.LState
	Thread  *lua.LState
}

type AddCharacter struct {
	Character  *entity.Character
	TargetMap  *entity.Map
	SpawnPoint uint8
	Init       bool
	Notes      []*internal.Note
}

type WarpCharacter struct {
	Character *entity.Character
	TargetMap *entity.Map
	Portal    uint8
	OnEnter   func(actor.Context)
	Ticket    *entity.MapRef
}

type HandoffCharacter struct {
	Character *entity.Character
	TargetMap *entity.Map
	Portal    uint8
	OnEnter   func(actor.Context)
	Ticket    *entity.MapRef
}

type TimerTick struct {
	HandlerName string
}

type ClearPartyByPartyID struct {
	CharacterID uint32
	PartyID     uint32
}

type SyncCharacterPartyState struct {
	CharacterID uint32
	PartyID     *uint32
}

type PartyMemberLeft struct {
	LeaverID uint32
}

type PartyDisband struct {
	FormerMemberIDs []uint32
}

type SaveMapCharacters struct{}

type SaveMapCharactersAck struct {
	MapID uint32
	Saved int
	Err   string
}

type ReceiveWhisper struct {
	CharacterID uint32
	Whisper     entity.Whisper
}

type ReceiveMessengerEvent struct {
	CharacterID uint32
	Event       entity.MessengerEvent
}

type SetFamily struct {
	CharacterID uint32
	Family      *internal.Family
}

type RefreshWedding struct {
	CharacterID uint32
	Event       string
}

type SpouseMoved struct {
	CharacterID uint32
	SpouseID    uint32
	MapID       uint32
	Reply       bool
}

type PartyMemberJoined struct {
	CharacterID     uint32
	ForChannel      int32
	PartyID         uint32
	JoinCharacterID uint32
	JoinName        string
	LeaderID        uint32
	Members         []response.PartyMemberStatus
}

type SetBuddyChannel struct {
	RecipientCharacterID uint32
	BuddyCharacterID     uint32
	Channel              int32
}

type SyncBuddyList struct {
	RecipientCharacterID uint32
	SyncAction           uint8
	Entries              []response.BuddyEntry
}

type InviteToGuild struct {
	CharacterID        uint32
	InviterCharacterID uint32
	GuildID            uint32
	InviterName        string
}

type InviteToAlliance struct {
	CharacterID        uint32
	TargetGuildID      uint32
	AllianceID         uint32
	InviterCharacterID uint32
	InviterGuildID     uint32
	InviterName        string
	AllianceName       string
	ExpiresAt          time.Time
}
