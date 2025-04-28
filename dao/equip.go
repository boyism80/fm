package dao

type Equip struct {
	UpgradeSlots  uint8  // 업그레이드 슬롯 수
	Level         uint8  // 장비 레벨
	Str           uint16 // Strength
	Dex           uint16 // Dexterity
	Int           uint16 // Intelligence
	Luk           uint16 // Luck
	Hp            uint16 // HP
	Mp            uint16 // MP
	Watk          uint16 // 공격력
	Matk          uint16 // 마법 공격력
	Wdef          uint16 // 물리 방어력
	Mdef          uint16 // 마법 방어력
	Acc           uint16 // 명중률
	Avoid         uint16 // 회피율
	Hands         uint16 // 손
	Speed         uint16 // 이동 속도
	Jump          uint16 // 점프
	IncSkill      uint16 // 스킬 증가
	BaseLevel     uint8  // 기본 레벨
	EquipLevel    uint8  // 장비 레벨
	ExpPercentage uint32 // 경험치 비율
}
