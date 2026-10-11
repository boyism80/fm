package wz

type ItemCore struct {
	ID              uint32
	Name            string
	Price           int
	Cash            bool
	Quest           bool
	SlotMax         uint16
	TradeAvailable  int
	TradeBlock      bool
	Only            bool
	AccountSharable bool
	MakerLevel      int
	MakerBonus      map[string]int16
}

func (model *ItemCore) setTradeFlags(info *node) {
	model.TradeBlock = info.Int("tradeBlock", 0) == 1
	model.Only = info.Int("only", 0) == 1
	model.AccountSharable = info.Int("accountSharable", 0) == 1
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

func (model *ItemCore) IsTradeBlock() bool {
	return model.TradeBlock
}

func (model *ItemCore) IsOnly() bool {
	return model.Only
}

func (model *ItemCore) IsAccountSharable() bool {
	return model.AccountSharable
}

func (model *ItemCore) GetMakerLevel() int { return model.MakerLevel }

func (model *ItemCore) GetMakerBonus() map[string]int16 { return model.MakerBonus }

func (model *ItemCore) IsConsumeOnPickup() bool {
	return false
}
