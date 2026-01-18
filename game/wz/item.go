package wz

type Item interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	IsCash() bool
	IsQuest() bool
	GetCapacity() uint16
	IsTradeAvailable() int
}

type BasicStats struct {
	Str uint16
	Dex uint16
	Int uint16
	Luk uint16
}

type RequiredStats struct {
	BasicStats
	Level uint8
	Class int
}

type AbilityStats struct {
	BasicStats
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

type ItemCore struct {
	ID             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
}

type Equipment struct {
	*ItemCore
	Required        RequiredStats
	Ability         AbilityStats
	TUC             uint8
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

type Pet struct {
	*ItemCore
}

type GeneralItem struct {
	*ItemCore
}

type CashItem struct {
	*ItemCore
}

type ActiveEffect struct {
	HP     int
	MP     int
	HPRate int // HP recovery rate (percentage, e.g., 100 = 100%)
	MPRate int // MP recovery rate (percentage, e.g., 100 = 100%)
}

type Consume struct {
	*ItemCore
	ActiveEffect ActiveEffect
}

type Installation struct {
	*ItemCore
}

type SpecialItem struct {
	*ItemCore
	ActiveEffect ActiveEffect
}

func (model *ItemCore) GetID() uint32 {
	return model.ID
}

func (model *ItemCore) GetName() string {
	return model.Name
}

func (model *ItemCore) GetPrice() int {
	return model.Price
}

func (model *ItemCore) IsCash() bool {
	return model.Cash
}

func (model *ItemCore) IsQuest() bool {
	return model.Quest
}

func (model *ItemCore) GetCapacity() uint16 {
	return max(1, model.SlotMax)
}

func (model *ItemCore) IsTradeAvailable() int {
	return model.TradeAvailable
}
