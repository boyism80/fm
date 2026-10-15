package game

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

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
			ClassId:     uint32(senior.Class),
		},
		Junior: &internal.FamilyMember{
			CharacterId: junior.GetID(),
			Name:        junior.GetName(),
			Level:       uint32(junior.GetLevel()),
			ClassId:     uint32(junior.Class),
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
