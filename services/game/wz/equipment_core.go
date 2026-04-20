package wz

type EquipmentCore struct {
	*ItemCore
	Required        RequiredStats
	Ability         AbilityStats
	EnchantChance   uint8
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

func (core *EquipmentCore) GetEnchantChance() uint8 { return core.EnchantChance }
func (core *EquipmentCore) GetRequired() RequiredStats {
	return core.Required
}
func (core *EquipmentCore) GetAbility() AbilityStats { return core.Ability }
