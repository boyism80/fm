package game

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

// MessengerAsync sends one messenger action for ch with its current look.
func (s *InternalSystem) MessengerAsync(ctx actor.Context, ch *Character, action internal.MessengerAction, roomID uint32, text string) *async.Promise[*internal.MessengerReply] {
	w := stream.NewStreamWriter(stream.LittleEndian)
	ch.BuildS2C().SerializeLook(w)
	req := &internal.MessengerRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		ChannelId:   uint32(s.gw.GetChannelID()),
		Action:      action,
		RoomId:      roomID,
		Text:        text,
		Look:        w.Bytes(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(
		func(c context.Context) (*internal.MessengerReply, error) {
			return s.client.Messenger(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) ClaimWhisperAsync(ctx actor.Context, ch *Character, requestID string) *async.Promise[*internal.CompleteWhisperReply] {
	req := &internal.CompleteWhisperRequest{
		WorldId:     s.gw.GetWorldID(),
		RequestId:   requestID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CompleteWhisperReply, error) {
			return s.client.CompleteWhisper(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) CompleteWhisperAsync(ctx actor.Context, ch *Character, requestID string, received bool) *async.Promise[*internal.CompleteWhisperReply] {
	req := &internal.CompleteWhisperRequest{
		WorldId:     s.gw.GetWorldID(),
		RequestId:   requestID,
		CharacterId: ch.GetID(),
		MapId:       ch.GetMap().TemplateID(),
		Complete:    true,
		Received:    received,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CompleteWhisperReply, error) {
			return s.client.CompleteWhisper(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) SendNoteAsync(ctx actor.Context, ch *Character, recipient string, message string) *async.Promise[*internal.NotesReply] {
	worldID := s.gw.GetWorldID()
	characterID := ch.GetID()
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return s.client.Notes(c, &internal.NotesRequest{
				WorldId:     worldID,
				CharacterId: characterID,
				Action:      internal.NotesAction_NOTES_ACTION_SEND,
				TargetName:  recipient,
				Message:     message,
				Gift:        uint32(pconst.NoteKindNormal),
			})
		},
		nil,
	)
}

// One update per character is in flight at a time; changes made meanwhile
// are sent as one update with the latest values when it ends.
func (s *InternalSystem) UpdateFamilyMemberAsync(ctx actor.Context, ch *Character) *async.Task {
	req := &internal.UpdateFamilyMemberRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Level:       uint32(ch.GetLevel()),
		ClassId:     uint32(ch.class),
	}
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	s.mu.Lock()
	if _, sending := s.familyUpdates[req.CharacterId]; sending {
		s.familyUpdates[req.CharacterId] = req
		s.mu.Unlock()
		return p
	}
	s.familyUpdates[req.CharacterId] = nil
	s.mu.Unlock()
	return s.sendFamilyUpdate(p, req)
}

func (s *InternalSystem) sendFamilyUpdate(p *async.Task, req *internal.UpdateFamilyMemberRequest) *async.Task {
	p.OnError(func(err error) {
		log.Printf("family update member: character=%d: %v", req.CharacterId, err)
	})
	p.Finally(func() {
		s.mu.Lock()
		next := s.familyUpdates[req.CharacterId]
		if next == nil {
			delete(s.familyUpdates, req.CharacterId)
			s.mu.Unlock()
			return
		}
		s.familyUpdates[req.CharacterId] = nil
		s.mu.Unlock()
		s.sendFamilyUpdate(async.NewTask(nil, core.InternalRPCPerStepTimeout), next)
	})
	p.ThenRPC(func(c context.Context) (*internal.UpdateFamilyMemberReply, error) {
		return s.client.UpdateFamilyMember(c, req)
	}, nil)
	return p
}

func (s *InternalSystem) GetFamilyAsync(ctx actor.Context, ch *Character, name string) *async.Promise[*internal.GetFamilyReply] {
	req := &internal.GetFamilyRequest{
		WorldId: s.gw.GetWorldID(),
		Name:    name,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.GetFamilyReply, error) {
			return s.client.GetFamily(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) JoinFamilyAsync(ctx actor.Context, senior *Character, junior *Character) *async.Promise[*internal.FamilyReply] {
	req := &internal.JoinFamilyRequest{
		WorldId: s.gw.GetWorldID(),
		Senior: &internal.FamilyMember{
			CharacterId: senior.GetID(),
			Name:        senior.GetName(),
			Level:       uint32(senior.GetLevel()),
			ClassId:     uint32(senior.class),
		},
		Junior: &internal.FamilyMember{
			CharacterId: junior.GetID(),
			Name:        junior.GetName(),
			Level:       uint32(junior.GetLevel()),
			ClassId:     uint32(junior.class),
		},
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return s.client.JoinFamily(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) SeparateFamilyAsync(ctx actor.Context, ch *Character, juniorID uint32) *async.Promise[*internal.FamilyReply] {
	req := &internal.SeparateFamilyRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		JuniorId:    juniorID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return s.client.SeparateFamily(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) SetFamilyPreceptAsync(ctx actor.Context, ch *Character, precept string) *async.Promise[*internal.FamilyReply] {
	req := &internal.SetFamilyPreceptRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Precept:     precept,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return s.client.SetFamilyPrecept(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) UseFamilyFeatureAsync(ctx actor.Context, ch *Character, feature internal.FamilyFeature, targetID uint32, requestID uint64) *async.Promise[*internal.FamilyReply] {
	req := &internal.UseFamilyFeatureRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Feature:     feature,
		TargetId:    targetID,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return s.client.UseFamilyFeature(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) FinishFamilyMoveAsync(ctx actor.Context, ch *Character, requestID uint64, complete bool) *async.Promise[*internal.FamilyReply] {
	req := &internal.FinishFamilyMoveRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		RequestId:   requestID,
		Complete:    complete,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return s.client.FinishFamilyMove(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) IncreaseBuddyCapacityAsync(ctx actor.Context, ch *Character, requestID uint64) *async.Promise[*internal.IncreaseBuddyCapacityReply] {
	req := &internal.IncreaseBuddyCapacityRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.IncreaseBuddyCapacityReply, error) {
			return s.client.IncreaseBuddyCapacity(c, req)
		},
		nil,
	)
}
