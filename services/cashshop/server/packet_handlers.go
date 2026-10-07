package server

import (
	"github.com/boyism80/fm/core"
)

func (cs *CashShopServer) registerPacketHandlers() {
	core.Bind[*CashShopServer, Pong](cs)
	core.Bind[*CashShopServer, Login](cs)
	core.Bind[*CashShopServer, Leave](cs)
	core.Bind[*CashShopServer, Refresh](cs)
	core.Bind[*CashShopServer, Operation](cs)
}
