package data

type baseStats struct {
	Str int
	Dex int
	Int int
	Luk int
}

type RequiredStats struct {
	baseStats
	Level int
	Class int
}

type AbilityStats struct {
	baseStats
	PAD       int // Physical attack damage
	MAD       int // Magical attack damage
	PDD       int // Physical defense damage
	MDD       int // Magical defense damage
	Speed     int
	Jump      int
	ACC       int
	EVA       int
	MaxHP     int
	MaxMP     int
	PVPDamage int
}

type ItemTemplate struct {
	Id              uint32
	Name            string
	Required        RequiredStats
	Ability         AbilityStats
	TUC             int // Total upgrade count
	Price           int
	AttackSpeed     int
	Cash            bool
	Quest           bool
	SlotMax         uint16
	EquipTradeBlock bool
	TradeAvailable  int
	Hide            bool
	RoyalSpecial    bool
	MasterSpecial   bool
}
