package wz

type EquipmentCore struct {
	*ItemCore
	Required        RequiredStats
	Ability         AbilityStats
	EnhanceChance   uint8
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
	Pets            []uint32
	TamingMob       uint32
}

func (core *EquipmentCore) GetEnhanceChance() uint8 { return core.EnhanceChance }
func (core *EquipmentCore) GetRequired() RequiredStats {
	return core.Required
}
func (core *EquipmentCore) GetAbility() AbilityStats { return core.Ability }
func (core *EquipmentCore) IsEquipTradeBlock() bool  { return core.EquipTradeBlock }
