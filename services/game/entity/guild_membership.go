package entity

import (
	"errors"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/constant"
)

var (
	ErrGuildAlreadyJoined  = errors.New("already in a guild")
	ErrGuildInviteBusy     = errors.New("guild invite already pending")
	ErrGuildInviteNotFound = errors.New("guild invite not found")
)

type GuildMembership struct {
	owner   *Character
	id      *uint32
	invites map[uint32]time.Time
}

func (g *GuildMembership) ID() (uint32, bool) {
	if g.id == nil {
		return 0, false
	}
	return *g.id, true
}

func (g *GuildMembership) SetID(id *uint32) {
	g.id = id
}

func (g *GuildMembership) Invite(guildID uint32, inviterName string) error {
	if g.id != nil {
		return ErrGuildAlreadyJoined
	}
	now := clock.Now()
	g.dropExpiredInvites(now)
	if len(g.invites) > 0 {
		return ErrGuildInviteBusy
	}

	g.invites[guildID] = now.Add(constant.GuildInviteDuration)
	g.owner.Listener.OnGuildInvite(g.owner, guildID, inviterName)
	return nil
}

func (g *GuildMembership) AcceptInvite(guildID uint32) error {
	if g.id != nil {
		return ErrGuildAlreadyJoined
	}
	g.dropExpiredInvites(clock.Now())
	if _, ok := g.invites[guildID]; !ok {
		return ErrGuildInviteNotFound
	}

	delete(g.invites, guildID)
	return nil
}

func (g *GuildMembership) ClearInvites() {
	clear(g.invites)
}

func (g *GuildMembership) dropExpiredInvites(now time.Time) {
	for id, expiresAt := range g.invites {
		if !now.Before(expiresAt) {
			delete(g.invites, id)
		}
	}
}
