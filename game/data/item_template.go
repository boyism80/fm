package data

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

type baseItemTemplate struct {
	Id             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
}

func (template *baseItemTemplate) GetID() uint32 {
	return template.Id
}

func (template *baseItemTemplate) GetName() string {
	return template.Name
}

func (template *baseItemTemplate) GetPrice() int {
	return template.Price
}

func (template *baseItemTemplate) IsCash() bool {
	return template.Cash
}

func (template *baseItemTemplate) IsQuest() bool {
	return template.Quest
}

func (template *baseItemTemplate) GetSlotMax() uint16 {
	return template.SlotMax
}

func (template *baseItemTemplate) IsTradeAvailable() int {
	return template.TradeAvailable
}

type EquipmentTemplate struct {
	*baseItemTemplate
	Required        RequiredStats
	Ability         AbilityStats
	TUC             uint8 // Total upgrade count
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

type PetTemplate struct {
	*baseItemTemplate
}

type GeneralItemTemplate struct {
	*baseItemTemplate
}

type CashItemTemplate struct {
	*baseItemTemplate
}

type ActiveEffect struct {
	HP int
	MP int
}

type ConsumeTemplate struct {
	*baseItemTemplate
	ActiveEffect ActiveEffect
}

type InstallationTemplate struct {
	*baseItemTemplate
}

type SpecialItemTemplate struct {
	*baseItemTemplate
	ActiveEffect ActiveEffect
}

type ItemTemplate interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	IsCash() bool
	IsQuest() bool
	GetSlotMax() uint16
	IsTradeAvailable() int
}
