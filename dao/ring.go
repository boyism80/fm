package dao

type Ring struct {
	RingId        uint64
	PartnerRingId uint64
	PartnerChrId  uint32
	ItemId        uint32
	PartnerName   string
	Equipped      bool
}

type RingCollection struct {
	Left  []*Ring
	Mid   []*Ring
	Right []*Ring
}
