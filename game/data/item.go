package data

type ItemSpec interface {
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

type ItemCoreSpec struct {
	Id             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
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

func (spec *ItemCoreSpec) GetID() uint32 {
	return spec.Id
}

func (spec *ItemCoreSpec) GetName() string {
	return spec.Name
}

func (spec *ItemCoreSpec) GetPrice() int {
	return spec.Price
}

func (spec *ItemCoreSpec) IsCash() bool {
	return spec.Cash
}

func (spec *ItemCoreSpec) IsQuest() bool {
	return spec.Quest
}

func (spec *ItemCoreSpec) GetCapacity() uint16 {
	return spec.SlotMax
}

func (spec *ItemCoreSpec) IsTradeAvailable() int {
	return spec.TradeAvailable
}
