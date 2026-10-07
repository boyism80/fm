package entity

import "time"

type MarriageStatus uint16

const (
	MarriageStatusEngaged MarriageStatus = 1
	MarriageStatusMarried MarriageStatus = 3
)

type Marriage struct {
	ID                 uint32
	GroomID            uint32
	BrideID            uint32
	GroomName          string
	BrideName          string
	GroomItemID        uint32
	BrideItemID        uint32
	Status             MarriageStatus
	TicketItemID       uint32
	GroomWished        bool
	BrideWished        bool
	GroomWishes        []string
	BrideWishes        []string
	DivorceRequesterID uint32
	DivorceRequestedAt time.Time
}

func (m *Marriage) PartnerID(characterID uint32) uint32 {
	if m.GroomID == characterID {
		return m.BrideID
	}
	return m.GroomID
}

func (m *Marriage) PartnerName(characterID uint32) string {
	if m.GroomID == characterID {
		return m.BrideName
	}
	return m.GroomName
}

func (m *Marriage) Wished(characterID uint32) bool {
	if m.GroomID == characterID {
		return m.GroomWished
	}
	return m.BrideWished
}

func (m *Marriage) Reserved() bool {
	return m.TicketItemID != 0 && m.GroomWished && m.BrideWished
}
