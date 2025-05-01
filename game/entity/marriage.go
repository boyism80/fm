package entity

type MarriageData struct {
	GroomId   uint32
	BrideId   uint32
	Status    uint16
	GroomName string
	BrideName string
}

var marriageManager *MarriageManager = &MarriageManager{}

type MarriageManager struct {
	marriages map[uint32]*MarriageData
}

func (m *MarriageManager) GetMarriage(marriageId uint32) *MarriageData {
	if m.marriages == nil {
		m.marriages = make(map[uint32]*MarriageData)
	}
	return m.marriages[marriageId]
}

// 전역 접근용
func GetMarriageManager() *MarriageManager {
	return marriageManager
}
