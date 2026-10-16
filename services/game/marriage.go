package game

import "time"

type MarriageStatus uint16

const (
	MarriageStatusEngaged MarriageStatus = 1
	MarriageStatusMarried MarriageStatus = 3
)

type Marriage struct {
	id                 uint32
	groomID            uint32
	brideID            uint32
	groomName          string
	brideName          string
	groomItemID        uint32
	brideItemID        uint32
	status             MarriageStatus
	ticketItemID       uint32
	groomWished        bool
	brideWished        bool
	groomWishes        []string
	brideWishes        []string
	divorceRequesterID uint32
	divorceRequestedAt time.Time
}

func (m *Marriage) PartnerID(characterID uint32) uint32 {
	if m.groomID == characterID {
		return m.brideID
	}
	return m.groomID
}

func (m *Marriage) PartnerName(characterID uint32) string {
	if m.groomID == characterID {
		return m.brideName
	}
	return m.groomName
}

func (m *Marriage) Wished(characterID uint32) bool {
	if m.groomID == characterID {
		return m.groomWished
	}
	return m.brideWished
}

func (m *Marriage) Reserved() bool {
	return m.ticketItemID != 0 && m.groomWished && m.brideWished
}
