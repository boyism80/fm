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
	core.Bind[*GameServer, HealOverTime](gs)
	core.Bind[*GameServer, CancelBuff](gs)
	core.Bind[*GameServer, UseItem](gs)
	core.Bind[*GameServer, UseInnerPortal](gs)
	core.Bind[*GameServer, PartySearchStop](gs)
	core.Bind[*GameServer, Unknown0C](gs)
	core.Bind[*GameServer, NpcShop](gs)
	core.Bind[*GameServer, MagicAttack](gs)
	core.Bind[*GameServer, ActiveSkill](gs)
	core.Bind[*GameServer, UseChair](gs)
	core.Bind[*GameServer, CancelChair](gs)
	core.Bind[*GameServer, ChangeKeymap](gs)
}
