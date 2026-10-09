package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type ExpeditionRegistry struct {
	mu     sync.Mutex
	gw     GameWorld
	boards map[string]*expeditionBoard
}

type expeditionBoard struct {
	spec         ExpeditionSpec
	recruiting   *Expedition
	battles      []*Expedition
	reservations []*expeditionReservation
}

type expeditionReservation struct {
	member *ExpeditionMember
	field  *Map
}

func NewExpeditionRegistry(gw GameWorld) *ExpeditionRegistry {
	return &ExpeditionRegistry{
		gw:     gw,
		boards: make(map[string]*expeditionBoard),
	}
}

func (b *expeditionBoard) reservationIndex(id uint32) int {
	for i, reservation := range b.reservations {
		if reservation.member.ID == id {
			return i
		}
	}
	return -1
}

func (r *ExpeditionRegistry) Register(name string, leader *Character, spec ExpeditionSpec) (*Expedition, error) {
	r.mu.Lock()
	board := r.boards[name]
	if board == nil {
		board = &expeditionBoard{}
		r.boards[name] = board
	}
	board.spec = spec
	if board.recruiting != nil {
		r.mu.Unlock()
		return nil, ErrExpeditionExists
	}
	if len(board.battles) >= spec.MaxBattles {
		r.mu.Unlock()
		return nil, ErrExpeditionBattling
	}

	e := &Expedition{
		Name:     name,
		Spec:     spec,
		Leader:   leader,
		BeginMap: leader.GetMap(),
		Deadline: time.Now().Add(time.Duration(spec.RecruitMs) * time.Millisecond),
		members:  []*ExpeditionMember{leader.ExpeditionMember()},
		registry: r,
	}
	if i := board.reservationIndex(leader.GetID()); i >= 0 {
		board.reservations = append(board.reservations[:i], board.reservations[i+1:]...)
	}
	e.expiry = time.AfterFunc(time.Duration(spec.RecruitMs)*time.Millisecond, func() {
		r.expire(e)
	})
	board.recruiting = e
	r.mu.Unlock()

	seconds := int32(spec.RecruitMs / 1000)
	notice := leader.GetName() + spec.Notice
	for _, obj := range e.BeginMap.GetAllPlayers() {
		ch, ok := obj.(*Character)
		if !ok {
			continue
		}
		ch.Listener.OnClock(ch, seconds)
		ch.Listener.OnMessage(ch, constant.MsgLightBlueText, notice)
	}
	return e, nil
}

func (r *ExpeditionRegistry) Find(name string) *Expedition {
	r.mu.Lock()
	defer r.mu.Unlock()
	board := r.boards[name]
	if board == nil {
		return nil
	}
	return board.recruiting
}

func (r *ExpeditionRegistry) Battles(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	board := r.boards[name]
	if board == nil {
		return 0
	}
	return len(board.battles)
}

func (r *ExpeditionRegistry) Reserve(name string, ch *Character) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	board := r.boards[name]
	if board == nil {
		board = &expeditionBoard{}
		r.boards[name] = board
	}

	if i := board.reservationIndex(ch.GetID()); i >= 0 {
		board.reservations = append(board.reservations[:i], board.reservations[i+1:]...)
		return false
	}
	board.reservations = append(board.reservations, &expeditionReservation{
		member: ch.ExpeditionMember(),
		field:  ch.GetMap(),
	})
	return true
}

func (r *ExpeditionRegistry) Reservations(name string) []ExpeditionMember {
	r.mu.Lock()
	defer r.mu.Unlock()
	board := r.boards[name]
	if board == nil {
		return nil
	}
	out := make([]ExpeditionMember, len(board.reservations))
	for i, reservation := range board.reservations {
		out[i] = *reservation.member
	}
	return out
}

func (r *ExpeditionRegistry) expire(e *Expedition) {
	r.mu.Lock()
	board := r.boards[e.Name]
	if board == nil || board.recruiting != e {
		r.mu.Unlock()
		return
	}
	board.recruiting = nil
	members := e.members
	r.mu.Unlock()

	e.stopClock()
	for _, member := range members {
		member.character.Listener.OnMessage(member.character, constant.MsgPinkText, constant.ExpeditionExpiredMessage)
	}
	r.promote(e.Name)
}

func (r *ExpeditionRegistry) disband(e *Expedition) error {
	r.mu.Lock()
	board := r.boards[e.Name]
	if board == nil || board.recruiting != e {
		r.mu.Unlock()
		return ErrExpeditionClosed
	}
	board.recruiting = nil
	e.expiry.Stop()
	members := e.members
	r.mu.Unlock()

	e.stopClock()
	for _, member := range members[1:] {
		member.character.Listener.OnMessage(member.character, constant.MsgPinkText, constant.ExpeditionDisbandedMessage)
	}
	r.promote(e.Name)
	return nil
}

func (r *ExpeditionRegistry) promote(name string) {
	r.mu.Lock()
	board := r.boards[name]
	if board == nil || board.recruiting != nil || len(board.battles) >= board.spec.MaxBattles || len(board.reservations) == 0 {
		r.mu.Unlock()
		return
	}
	reservation := board.reservations[0]
	board.reservations = board.reservations[1:]
	spec := board.spec
	r.mu.Unlock()

	r.gw.GetMapSystem().Call(reservation.field, func(actor.Context) {
		leader := reservation.field.GetPlayer(reservation.member.ID)
		if leader == nil {
			notice := fmt.Sprintf(constant.ExpeditionSkippedMessage, reservation.member.Name)
			for _, obj := range reservation.field.GetAllPlayers() {
				if ch, ok := obj.(*Character); ok {
					ch.Listener.OnMessage(ch, constant.MsgLightBlueText, notice)
				}
			}
			r.promote(name)
			return
		}

		if _, err := r.Register(name, leader, spec); err != nil {
			r.promote(name)
		}
	})
}

func (r *ExpeditionRegistry) LeaveChannel(ch *Character) {
	id := ch.GetID()
	disbanded := make([]*Expedition, 0)
	left := make([]*Expedition, 0)

	r.mu.Lock()
	for _, board := range r.boards {
		if i := board.reservationIndex(id); i >= 0 {
			board.reservations = append(board.reservations[:i], board.reservations[i+1:]...)
		}

		e := board.recruiting
		if e == nil {
			continue
		}
		if e.Leader.GetID() == id {
			disbanded = append(disbanded, e)
			continue
		}
		if i := e.memberIndex(id); i >= 0 {
			e.members = append(e.members[:i], e.members[i+1:]...)
			left = append(left, e)
		}
		if i := e.bannedIndex(id); i >= 0 {
			e.banned[i].character = nil
		}
	}
	r.mu.Unlock()

	for _, e := range disbanded {
		_ = r.disband(e)
	}
	for _, e := range left {
		e.Leader.Listener.OnMessage(e.Leader, constant.MsgPinkText, fmt.Sprintf(constant.ExpeditionLeftMessage, ch.GetName()))
	}
}

func (r *ExpeditionRegistry) EndBattle(sm *StateMachine) {
	r.mu.Lock()
	name := ""
	for boardName, board := range r.boards {
		for i, e := range board.battles {
			if e.sm == sm {
				board.battles = append(board.battles[:i], board.battles[i+1:]...)
				name = boardName
				break
			}
		}
	}
	r.mu.Unlock()
	if name == "" {
		return
	}

	r.promote(name)
}

func (r *ExpeditionRegistry) SyncClock(ch *Character, m *Map) {
	r.mu.Lock()
	var recruiting *Expedition
	for _, board := range r.boards {
		if board.recruiting != nil && board.recruiting.BeginMap == m {
			recruiting = board.recruiting
			break
		}
	}
	r.mu.Unlock()
	if recruiting == nil {
		return
	}

	left := recruiting.TimeLeft()
	if left <= 0 {
		return
	}
	ch.Listener.OnClock(ch, int32(left/1000))
}
