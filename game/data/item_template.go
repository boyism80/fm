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
	TUC             int // Total upgrade count
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

type ActiveEffect struct {
	HP int
	MP int
}

type ConsumeTemplate struct {
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
