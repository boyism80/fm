package game

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
	lua "github.com/yuin/gopher-lua"
)

type RequestSpawnDoor struct {
	ReplyTo        *actor.PID
	TargetMap      *Map
	Key            DoorKey
	CharacterID    uint32
	OwnerID        uint32
	SkillID        constant.SkillID
	Field          DoorEndpoint
	PartyOwnerSlot int
	PartyID        *uint32
}

type ResponseSpawnDoor struct {
	Ok          bool
	Key         DoorKey
	CharacterID uint32
	OwnerID     uint32
	SkillID     constant.SkillID
	Return      DoorEndpoint
	Field       DoorEndpoint
	PartyID     *uint32
}

type DespawnDoor struct {
	Map               *Map
	Key               DoorKey
	Animated          bool
	NotifyCounterpart bool
}

type SetDoorPartyID struct {
	Map     *Map
	Key     DoorKey
	PartyID *uint32
}

type ResumeLua struct {
	Root   *lua.LState
	Thread *lua.LState
	Args   []lua.LValue
}

type MapCallFunc func(ctx actor.Context, maps []*Map) []lua.LValue

type MapCallAsyncFunc func(ctx actor.Context, maps []*Map) *async.Promise[[]lua.LValue]

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
	Character  *Character
	TargetMap  *Map
	SpawnPoint uint8
	Init       bool
	Notes      []Note
}

type WarpCharacter struct {
	Character *Character
	TargetMap *Map
	Portal    uint8
	OnEnter   func(actor.Context)
	Ticket    *MapRef
}

type HandoffCharacter struct {
	Character *Character
	TargetMap *Map
	Portal    uint8
	OnEnter   func(actor.Context)
	Ticket    *MapRef
}

type ClearPartyByPartyID struct {
	CharacterID uint32
	PartyID     uint32
}

type SyncCharacterPartyState struct {
	CharacterID uint32
	PartyID     *uint32
	UpdatedAt   time.Time
	Clear       bool
}

type PartyMemberLeft struct {
	LeaverID uint32
}

type PartyDisband struct {
	CharacterID     uint32
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
	Whisper     Whisper
}

type ReceiveMessengerEvent struct {
	CharacterID uint32
	Event       MessengerEvent
}

type SetFamily struct {
	CharacterID uint32
	Family      FamilyTree
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
