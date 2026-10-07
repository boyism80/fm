package constant

const (
	RingBoxFirst           uint32 = 2240004
	RingBoxLast            uint32 = 2240015
	EngagementRingFirst    uint32 = 4210000
	EngagementRingLast     uint32 = 4210011
	WeddingRingFirst       uint32 = 1112300
	WeddingRingLast        uint32 = 1112311
	WeddingInvitationFirst uint32 = 4211000
	WeddingInvitationLast  uint32 = 4211002
	WeddingTicketFirst     uint32 = 5251004
)

type WeddingTicket struct {
	Invitation      uint32
	InvitationCount uint16
	InvitedItem     uint32
}

var WeddingTickets = map[uint32]WeddingTicket{
	5251004: {Invitation: 4211000, InvitationCount: 5, InvitedItem: 4212000},
	5251005: {Invitation: 4211001, InvitationCount: 15, InvitedItem: 4212001},
	5251006: {Invitation: 4211002, InvitationCount: 30, InvitedItem: 4212002},
}

func IsEngagementRing(itemID uint32) bool {
	return itemID >= EngagementRingFirst && itemID <= EngagementRingLast
}

func IsWeddingRing(itemID uint32) bool {
	return itemID >= WeddingRingFirst && itemID <= WeddingRingLast
}
