package wz

type ItemCore struct {
	ID             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
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

func (model *ItemCore) IsConsumeOnPickup() bool {
	return false
}
