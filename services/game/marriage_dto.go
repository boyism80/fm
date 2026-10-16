package game

import "github.com/boyism80/fm/protocol/dto"

func (m *Marriage) BuildS2C() *dto.Marriage {
	return &dto.Marriage{
		ID:          m.id,
		GroomID:     m.groomID,
		BrideID:     m.brideID,
		Status:      uint16(m.status),
		GroomItemID: m.groomItemID,
		BrideItemID: m.brideItemID,
		GroomName:   m.groomName,
		BrideName:   m.brideName,
	}
}
