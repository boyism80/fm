package data

type ItemSpec interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	IsCash() bool
	IsQuest() bool
	GetSlotMax() uint16
	IsTradeAvailable() int
}

type baseStats struct {
	Str uint16
	Dex uint16
	Int uint16
	Luk uint16
}

type RequiredStats struct {
	baseStats
	Level uint8
	Class int
}

type AbilityStats struct {
	baseStats
	PAD       uint16 // Physical attack damage
	MAD       uint16 // Magical attack damage
	PDD       uint16 // Physical defense damage
	MDD       uint16 // Magical defense damage
	Speed     uint16
	Jump      uint16
	ACC       uint16
	EVA       int
	MaxHP     uint16
	MaxMP     uint16
	PVPDamage int

	Avoid uint16
	Hands uint16
}

type ItemCoreSpec struct {
	Id             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
}

func (template *ItemCoreSpec) GetID() uint32 {
	return template.Id
}

func (template *ItemCoreSpec) GetName() string {
	return template.Name
}

func (template *ItemCoreSpec) GetPrice() int {
	return template.Price
}

func (template *ItemCoreSpec) IsCash() bool {
	return template.Cash
}

func (template *ItemCoreSpec) IsQuest() bool {
	return template.Quest
}

func (template *ItemCoreSpec) GetSlotMax() uint16 {
	return template.SlotMax
}

func (template *ItemCoreSpec) IsTradeAvailable() int {
	return template.TradeAvailable
}

type EquipmentSpec struct {
	*ItemCoreSpec
	Required        RequiredStats
	Ability         AbilityStats
	TUC             uint8 // Total upgrade count
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

type PetSpec struct {
	*ItemCoreSpec
}

type GeneralItemSpec struct {
	*ItemCoreSpec
}

type CashItemSpec struct {
	*ItemCoreSpec
}

type ActiveEffect struct {
	HP int
	MP int
}

type ConsumeSpec struct {
	*ItemCoreSpec
	ActiveEffect ActiveEffect
}

type InstallationSpec struct {
	*ItemCoreSpec
}

type SpecialItemSpec struct {
	*ItemCoreSpec
	ActiveEffect ActiveEffect
}
