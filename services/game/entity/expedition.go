package entity

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

var (
	ErrExpeditionClosed    = errors.New("expedition closed")
	ErrExpeditionExists    = errors.New("expedition exists")
	ErrExpeditionBattling  = errors.New("expedition battles full")
	ErrExpeditionFull      = errors.New("expedition full")
	ErrExpeditionJoined    = errors.New("expedition joined")
	ErrExpeditionReserved  = errors.New("expedition reserved")
	ErrExpeditionBanned    = errors.New("expedition banned")
	ErrExpeditionNotMember = errors.New("expedition not member")
	ErrExpeditionNotLeader = errors.New("expedition not leader")
	ErrExpeditionAway      = errors.New("expedition member away")
	ErrExpeditionTooFew    = errors.New("expedition too few")
)

type ExpeditionSpec struct {
	Group      string
	MinMembers int
	MaxMembers int
	MaxBattles int
	RecruitMs  int64
	Notice     string
}

type ExpeditionRole int

const (
	ExpeditionRoleNone ExpeditionRole = iota
	ExpeditionRoleLeader
	ExpeditionRoleMember
	ExpeditionRoleBanned
)

type ExpeditionMember struct {
	ID        uint32
	Name      string
	Class     uint16
	character *Character
}

type Expedition struct {
	Name     string
	Spec     ExpeditionSpec
	Leader   *Character
	BeginMap *Map
	Deadline time.Time
	members  []*ExpeditionMember
	banned   []*ExpeditionMember
	sm       *StateMachine
	expiry   *time.Timer
	registry *ExpeditionRegistry
}

func (ch *Character) ExpeditionMember() *ExpeditionMember {
	return &ExpeditionMember{
		ID:        ch.GetID(),
		Name:      ch.GetName(),
		Class:     ch.Class,
		character: ch,
	}
}

func (e *Expedition) openLocked() bool {
	board := e.registry.boards[e.Name]
	return board != nil && board.recruiting == e
}

func (e *Expedition) memberIndex(id uint32) int {
	for i, member := range e.members {
		if member.ID == id {
			return i
		}
	}
	return -1
}

func (e *Expedition) bannedIndex(id uint32) int {
	for i, member := range e.banned {
		if member.ID == id {
			return i
		}
	}
	return -1
}

func (e *Expedition) Join(ch *Character) error {
	err := e.addMember(ch)
	if err != nil {
		return err
	}

	e.Leader.Listener.OnMessage(e.Leader, constant.MsgPinkText, fmt.Sprintf(constant.ExpeditionJoinedMessage, ch.GetName()))
	return nil
}

func (e *Expedition) addMember(ch *Character) error {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return ErrExpeditionClosed
	}
	if e.memberIndex(ch.GetID()) >= 0 {
		return ErrExpeditionJoined
	}
	if i := e.bannedIndex(ch.GetID()); i >= 0 {
		e.banned[i].character = ch
		return ErrExpeditionBanned
	}
	if e.registry.boards[e.Name].reservationIndex(ch.GetID()) >= 0 {
		return ErrExpeditionReserved
	}
	if len(e.members) >= e.Spec.MaxMembers {
		return ErrExpeditionFull
	}

	e.members = append(e.members, ch.ExpeditionMember())
	return nil
}

func (e *Expedition) Leave(ch *Character) error {
	if ch.GetID() == e.Leader.GetID() {
		return e.registry.disband(e)
	}

	err := e.removeMember(ch.GetID())
	if err != nil {
		return err
	}

	e.Leader.Listener.OnMessage(e.Leader, constant.MsgPinkText, fmt.Sprintf(constant.ExpeditionLeftMessage, ch.GetName()))
	return nil
}

func (e *Expedition) removeMember(id uint32) error {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return ErrExpeditionClosed
	}
	i := e.memberIndex(id)
	if i < 0 {
		return ErrExpeditionNotMember
	}

	e.members = append(e.members[:i], e.members[i+1:]...)
	return nil
}

func (e *Expedition) Kick(leader *Character, memberID uint32) error {
	kicked, err := e.ban(leader, memberID)
	if err != nil {
		return err
	}

	kicked.Listener.OnMessage(kicked, constant.MsgPinkText, constant.ExpeditionKickedMessage)
	return nil
}

func (e *Expedition) ban(leader *Character, memberID uint32) (*Character, error) {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return nil, ErrExpeditionClosed
	}
	if leader.GetID() != e.Leader.GetID() {
		return nil, ErrExpeditionNotLeader
	}
	if memberID == e.Leader.GetID() {
		return nil, ErrExpeditionNotMember
	}
	i := e.memberIndex(memberID)
	if i < 0 {
		return nil, ErrExpeditionNotMember
	}

	member := e.members[i]
	e.members = append(e.members[:i], e.members[i+1:]...)
	e.banned = append(e.banned, member)
	return member.character, nil
}

func (e *Expedition) Allow(leader *Character, memberID uint32) error {
	allowed, err := e.unban(leader, memberID)
	if err != nil {
		return err
	}

	allowed.Listener.OnMessage(allowed, constant.MsgPinkText, constant.ExpeditionAllowedMessage)
	return nil
}

func (e *Expedition) unban(leader *Character, memberID uint32) (*Character, error) {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return nil, ErrExpeditionClosed
	}
	if leader.GetID() != e.Leader.GetID() {
		return nil, ErrExpeditionNotLeader
	}
	i := e.bannedIndex(memberID)
	if i < 0 {
		return nil, ErrExpeditionNotMember
	}
	member := e.banned[i]
	if member.character == nil {
		return nil, ErrExpeditionAway
	}
	if len(e.members) >= e.Spec.MaxMembers {
		return nil, ErrExpeditionFull
	}

	e.banned = append(e.banned[:i], e.banned[i+1:]...)
	e.members = append(e.members, member)
	return member.character, nil
}

func (e *Expedition) Role(ch *Character) ExpeditionRole {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return ExpeditionRoleNone
	}

	switch {
	case ch.GetID() == e.Leader.GetID():
		return ExpeditionRoleLeader
	case e.memberIndex(ch.GetID()) >= 0:
		return ExpeditionRoleMember
	case e.bannedIndex(ch.GetID()) >= 0:
		return ExpeditionRoleBanned
	default:
		return ExpeditionRoleNone
	}
}

func (e *Expedition) Members() []ExpeditionMember {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	out := make([]ExpeditionMember, len(e.members))
	for i, member := range e.members {
		out[i] = *member
	}
	return out
}

func (e *Expedition) Banned() []ExpeditionMember {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	out := make([]ExpeditionMember, len(e.banned))
	for i, member := range e.banned {
		out[i] = *member
	}
	return out
}

func (e *Expedition) TimeLeft() int64 {
	left := time.Until(e.Deadline).Milliseconds()
	if left < 0 {
		return 0
	}
	return left
}

func (e *Expedition) Start(leader *Character) (*StateMachine, []*Character, error) {
	sm, entrants, skipped, err := e.startBattle(leader)
	if err != nil {
		return nil, nil, err
	}

	e.stopClock()
	for _, ch := range entrants {
		sm.EnterPlayer(ch)
	}
	sm.Start()
	return sm, skipped, nil
}

func (e *Expedition) stopClock() {
	e.registry.gw.GetMapSystem().Call(e.BeginMap, func(actor.Context) {
		for _, obj := range e.BeginMap.GetAllPlayers() {
			if ch, ok := obj.(*Character); ok {
				ch.Listener.OnStopClock(ch)
			}
		}
	})
}

func (e *Expedition) startBattle(leader *Character) (*StateMachine, []*Character, []*Character, error) {
	e.registry.mu.Lock()
	defer e.registry.mu.Unlock()
	if !e.openLocked() {
		return nil, nil, nil, ErrExpeditionClosed
	}
	if leader.GetID() != e.Leader.GetID() {
		return nil, nil, nil, ErrExpeditionNotLeader
	}
	board := e.registry.boards[e.Name]
	if len(board.battles) >= e.Spec.MaxBattles {
		return nil, nil, nil, ErrExpeditionBattling
	}

	field := leader.GetMap()
	entrants := make([]*Character, 0, len(e.members))
	skipped := make([]*Character, 0)
	for _, member := range e.members {
		if member.character.GetMap() != field {
			continue
		}
		if member.character.StateMachine() != nil {
			skipped = append(skipped, member.character)
			continue
		}
		entrants = append(entrants, member.character)
	}
	exempt := leader.ActsAsGM()
	if len(entrants) < e.Spec.MinMembers && !exempt {
		return nil, nil, nil, ErrExpeditionTooFew
	}

	group := e.registry.gw.GetStateMachineRegistry().Get(e.Spec.Group)
	if group == nil {
		return nil, nil, nil, fmt.Errorf("state machine group %s not found", e.Spec.Group)
	}
	sm, err := group.Create(strconv.FormatUint(uint64(leader.GetID()), 10), CreateOpts{
		Leader: leader,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	e.sm = sm
	e.expiry.Stop()
	board.recruiting = nil
	board.battles = append(board.battles, e)
	return sm, entrants, skipped, nil
}
