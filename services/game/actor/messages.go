package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type RequestSpawnDoor struct {
	ReplyTo        *actor.PID
	TargetMapID    uint32
	CharacterID    uint32
	OwnerID        uint32
	SkillID        constant.SkillID
	Field          entity.DoorEndpoint
	PartyOwnerSlot int
	PartyID        *uint32
}

type ResponseSpawnDoor struct {
	Ok          bool
	CharacterID uint32
	OwnerID     uint32
	SkillID     constant.SkillID
	Return      entity.DoorEndpoint
	Field       entity.DoorEndpoint
	PartyID     *uint32
}

type RemoveDoor struct {
	MapID   uint32
	OwnerID uint32
	SkillID uint32
}

type ResumeLua struct {
	Root   *lua.LState
	Thread *lua.LState
	Args   []lua.LValue
}

type MapCallFunc func(ctx actor.Context, a *MapActor) []lua.LValue

type MapCallAsyncFunc func(ctx actor.Context, a *MapActor) *async.Promise

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

type MapCallAck struct {
	Root   *lua.LState
	Thread *lua.LState
	Values []lua.LValue
}

type AddCharacter struct {
	Character  *entity.Character
	TargetMap  *entity.Map
	SpawnPoint uint8
	Init       bool
}

type RemoveCharacter struct {
	CharacterID uint32
}

type WarpCharacter struct {
	Character *entity.Character
	TargetMap *entity.Map
	Portal    uint8
}

type HandoffCharacter struct {
	Character *entity.Character
	TargetMap *entity.Map
	Portal    uint8
}

type AttachStateMachine struct {
	StateMachine *entity.StateMachine
	ReplyTo      *actor.PID
	MapID        uint32
}

type AttachStateMachineAck struct {
	MapID uint32
	OK    bool
}

type DetachStateMachine struct {
	StateMachine *entity.StateMachine
	ReplyTo      *actor.PID
	MapID        uint32
}

type DetachStateMachineAck struct {
	MapID uint32
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

type DeliverGuildNewMember struct {
	CharacterID uint32
	GuildID     uint32
	Member      dto.GuildMemberStatus
}

type DeliverGuildLeaveSelf struct {
	CharacterID uint32
}

type DeliverGuildExpelSelf struct {
	CharacterID uint32
	GuildID     uint32
}

type DeliverGuildMemberLeft struct {
	CharacterID uint32
	GuildID     uint32
	TargetID    uint32
	TargetName  string
	WasExpelled bool
}

type DeliverGuildRankTitleChange struct {
	CharacterID uint32
	GuildID     uint32
	RankTitles  [5]string
}

type DeliverGuildMemberRankChange struct {
	CharacterID uint32
	GuildID     uint32
	TargetID    uint32
	GuildRank   uint8
}

type DeliverGuildEmblemChange struct {
	CharacterID uint32
	GuildID     uint32
	LogoBG      uint16
	LogoBGColor uint8
	Logo        uint16
	LogoColor   uint8
}

type DeliverGuildNoticeChange struct {
	CharacterID uint32
	GuildID     uint32
	Notice      string
}

type DeliverGuildCapacityChange struct {
	CharacterID uint32
	GuildID     uint32
	Capacity    uint8
}

type DeliverGuildMemberOnlineChange struct {
	CharacterID uint32
	GuildID     uint32
	SubjectID   uint32
	Online      bool
}

type DeliverGuildMemberFieldsChange struct {
	CharacterID uint32
	GuildID     uint32
	SubjectID   uint32
	Level       uint32
	ClassID     uint32
}

type DeliverGuildDisbandSelf struct {
	CharacterID uint32
	GuildID     uint32
}

type DeliverGuildMessage struct {
	CharacterID uint32
	Code        pconst.GuildResponseCode
}

type DeliverAllianceMemberOnlineChange struct {
	CharacterID uint32
	AllianceID  uint32
	GuildID     uint32
	SubjectID   uint32
	Online      bool
}

type DeliverAllianceMemberFieldsChange struct {
	CharacterID uint32
	AllianceID  uint32
	GuildID     uint32
	SubjectID   uint32
	Level       uint32
	ClassID     uint32
}

type DeliverAllianceCreate struct {
	CharacterID      uint32
	Info             *dto.AllianceInfo
	Guilds           []*dto.GuildInfo
	MembershipGuilds []dto.AllianceMembershipChangeGuild
}

type DeliverAllianceDisband struct {
	CharacterID uint32
	AllianceID  uint32
}

type DeliverAllianceInfoBroadcast struct {
	CharacterID uint32
	Info        *dto.AllianceInfo
}

type DeliverAllianceNoticeChanged struct {
	CharacterID uint32
	Info        *dto.AllianceInfo
}

type DeliverAllianceLeaderChanged struct {
	CharacterID uint32
	AllianceID  uint32
	OldLeaderID uint32
	NewLeaderID uint32
	Info        *dto.AllianceInfo
	Guilds      []*dto.GuildInfo
}

type DeliverAllianceMemberRankChanged struct {
	CharacterID uint32
	Info        *dto.AllianceInfo
	Guilds      []*dto.GuildInfo
}

type DeliverAllianceGuildAdded struct {
	CharacterID     uint32
	Info            *dto.AllianceInfo
	Guilds          []*dto.GuildInfo
	NewGuildID      uint32
	AddedGuild      *dto.GuildInfo
	Members         []dto.AllianceGuildMemberRank
	Joining         bool
	MembershipGuild dto.AllianceMembershipChangeGuild
	HasMembership   bool
}

type DeliverAllianceGuildLeft struct {
	CharacterID    uint32
	Info           *dto.AllianceInfo
	RemovedGuildID uint32
	RemovedGuild   *dto.GuildInfo
	RemovedMembers []dto.AllianceGuildMemberRank
	Expelled       bool
	Leaving        bool
}

type DeliverAllianceInvite struct {
	CharacterID        uint32
	TargetGuildID      uint32
	AllianceID         uint32
	InviterCharacterID uint32
	InviterGuildID     uint32
	InviterName        string
	AllianceName       string
	ExpiresAt          time.Time
}
