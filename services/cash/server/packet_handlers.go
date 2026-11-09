package server

func (cs *CashShopServer) registerPacketHandlers() {
	cs.Bind[Pong]()
	cs.Bind[Login]()
	cs.Bind[Leave]()
	cs.Bind[Refresh]()
	cs.Bind[Operation]()
	cs.Bind[Coupon]()
}
