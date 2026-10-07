package entity

import "github.com/boyism80/fm/protocol/dto"

func (m *Marriage) ToDTO() *dto.Marriage {
	return &dto.Marriage{
		ID:          m.ID,
		GroomID:     m.GroomID,
		BrideID:     m.BrideID,
		Status:      uint16(m.Status),
		GroomItemID: m.GroomItemID,
		BrideItemID: m.BrideItemID,
		GroomName:   m.GroomName,
		BrideName:   m.BrideName,
	}
}
