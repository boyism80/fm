package entity

type Ring struct {
	RingId       uint64
	PartnerId    uint64
	RingUniqueId uint64
	PartnerChrId uint32
	ItemId       uint32
	PartnerName  string
	Equipped     bool
}

// TODO: load cash couple and friendship rings into Left and Mid (#258)
type RingContainer struct {
	Left []*Ring
	Mid  []*Ring
}
