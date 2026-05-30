package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
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
	Args   []lua.LValue
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

type SyncParty struct {
	Party *entity.Party
}

type ClearPartyByPartyID struct {
	PartyID uint32
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

type DeliverPartyInvite struct {
	CharacterID uint32
	PartyID     uint32
	InviterName string
	PartySearch bool
}

type DeliverMultiChat struct {
	CharacterID uint32
	Mode        pconst.MultiChatMode
	SenderName  string
	Message     string
}

type DeliverPartyStatusMessage struct {
	CharacterID uint32
	Code        pconst.PartyStatusCode
	Name        string
}

type DeliverPartyUpdateJoin struct {
	CharacterID     uint32
	ForChannel      int32
	PartyID         uint32
	JoinCharacterID uint32
	JoinName        string
	LeaderID        uint32
	Members         []response.PartyMemberStatus
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

type DeliverBuddyChannelUpdate struct {
	RecipientCharacterID uint32
	BuddyCharacterID     uint32
	Channel              int32
}

type DeliverBuddyListUpdate struct {
	RecipientCharacterID uint32
	SyncAction           uint8
	Entries              []response.BuddyEntry
}

type DeliverBuddyAddRequest struct {
	RecipientCharacterID uint32
	FromCharacterID      uint32
	FromName             string
}

type DeliverMessage struct {
	CharacterID uint32
	MessageType constant.ServerMessageType
	Message     string
}

type DeliverGuildInvite struct {
	CharacterID        uint32
	InviterCharacterID uint32
	GuildID            uint32
	InviterName        string
}

type DeliverGuildMessage struct {
	CharacterID uint32
	Code        pconst.GuildResponseCode
}
