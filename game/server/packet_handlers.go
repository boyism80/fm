package server

import (
	"github.com/boyism80/fm/core"
)

func (gs *GameServer) registerPacketHandlers() {
	core.Bind[*GameServer, Pong](gs)
	core.Bind[*GameServer, LoginGame](gs)
	core.Bind[*GameServer, MovePlayer](gs)
	core.Bind[*GameServer, NormalChat](gs)
	core.Bind[*GameServer, Attack](gs)
	core.Bind[*GameServer, ItemLoot](gs)
	core.Bind[*GameServer, DropMeso](gs)
	core.Bind[*GameServer, NpcControl](gs)
	core.Bind[*GameServer, Dialog](gs)
	core.Bind[*GameServer, NpcClick](gs)
	core.Bind[*GameServer, MoveMob](gs)
	core.Bind[*GameServer, Damaged](gs)
	core.Bind[*GameServer, Warp](gs)
	core.Bind[*GameServer, MoveItem](gs)
	core.Bind[*GameServer, SortInventory](gs)
	core.Bind[*GameServer, DistributeAP](gs)
	core.Bind[*GameServer, AutoAssignAP](gs)
	core.Bind[*GameServer, DistributeSP](gs)
}
