package game

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/rpc"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type InternalSystem struct {
	gw     GameWorld
	client internal.InternalClient
}

func NewInternalSystem(gw GameWorld, conn *rpc.InternalConnection) *InternalSystem {
	return &InternalSystem{gw: gw, client: conn.Client}
}

func (s *InternalSystem) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
	reply, err := s.client.CloseChannelShops(ctx, &internal.CloseChannelShopsRequest{
		WorldId:   s.gw.GetWorldID(),
		ChannelId: s.gw.GetChannelID(),
	})
	cancel()
	if err != nil {
		log.Printf("close channel shops: %v", err)
	} else {
		log.Printf("Closed %d shops into the store bank", reply.GetCount())
	}
	return nil
}

func (s *InternalSystem) Stop() {}
