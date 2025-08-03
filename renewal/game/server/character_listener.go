package server

import (
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/resp"
)

// GameCharacterListener implements CharacterListener for game server
type GameCharacterListener struct {
	gameServer *GameServer
	character  *entity.Character
}

// NewGameCharacterListener creates a new GameCharacterListener instance
func NewGameCharacterListener(gameServer *GameServer, character *entity.Character) *GameCharacterListener {
	return &GameCharacterListener{
		gameServer: gameServer,
		character:  character,
	}
}

// OnDialog handles dialog events
func (l *GameCharacterListener) OnDialog(npc uint32, message string, prev bool, next bool) {
	// TODO: Implement dialog handling
}

// OnDialogYesNo handles yes/no dialog events
func (l *GameCharacterListener) OnDialogYesNo(npc uint32, message string, prev bool, next bool) {
	// TODO: Implement yes/no dialog handling
}

// OnDialogAccept handles dialog accept events
func (l *GameCharacterListener) OnDialogAccept(npc uint32, message string, enableEscape bool) {
	// TODO: Implement dialog accept handling
}

// OnDialogList handles dialog list events
func (l *GameCharacterListener) OnDialogList(npc uint32, message string, selections []string) {
	// TODO: Implement dialog list handling
}

// OnDialogInput handles dialog input events
func (l *GameCharacterListener) OnDialogInput(npc uint32, message string) {
	// TODO: Implement dialog input handling
}

// OnChat handles chat events by broadcasting chat packet to all players on the map
func (l *GameCharacterListener) OnChat(message string, highlight bool, dontRecordHistory bool) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(l.character.GetMap())
	if mapInstance == nil {
		return
	}

	// Create chat packet
	chatPacket := &resp.NormalChat{
		CharacterId:       l.character.GetID(),
		Highlight:         highlight,
		Message:           message,
		DontRecordHistory: dontRecordHistory,
	}

	// Broadcast to all players on the map (including sender)
	mapInstance.BroadcastToAllPlayers(chatPacket, types.SEND_POLICY_ENCRYPT)
}

// OnAttack handles attack events by broadcasting attack packet to all players on the map
func (l *GameCharacterListener) OnAttack(attackInfo protocol.AttackInfo) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(l.character.GetMap())
	if mapInstance == nil {
		return
	}

	// Create attack packet
	attackPacket := &resp.Attack{
		AttackInfo:  attackInfo,
		CharacterId: l.character.GetID(),
		SkillLevel:  0, // TODO: Get actual skill level
	}

	// Broadcast to all players on the map (including sender)
	mapInstance.BroadcastToAllPlayers(attackPacket, types.SEND_POLICY_ENCRYPT)
}

// OnMesoChanged handles meso change events
func (l *GameCharacterListener) OnMesoChanged(meso int32) {
	// TODO: Implement meso change handling
}

// OnMessage handles general message events
func (l *GameCharacterListener) OnMessage(messageType constant.ServerMessageType, message string) {
	// Send notice packet to the character
	l.character.Send(&resp.Notice{
		Type:    messageType,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnExpGain handles experience gain events by sending exp gain packet to the character
func (l *GameCharacterListener) OnExpGain(exp uint32) {
	// Send exp gain packet to the character (following old server pattern)
	expPacket := &resp.GainExp{
		Gain:  exp,
		White: false,
	}

	l.character.Send(expPacket, types.SEND_POLICY_ENCRYPT)
	l.character.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_EXP: int32(l.character.Exp),
		},
	}, types.SEND_POLICY_ENCRYPT)
}
