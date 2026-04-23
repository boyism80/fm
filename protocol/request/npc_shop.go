package request

import "github.com/boyism80/fm/stream"

type ShopTransaction interface {
	GetMode() byte
	Deserialize(reader *stream.StreamReader)
}

type NpcShop struct {
	Mode        byte
	Transaction ShopTransaction
}

type BuyTransaction struct {
	Unknown  int16
	ItemID   uint32
	Quantity uint16
}

func (t *BuyTransaction) GetMode() byte {
	return 0
}

func (t *BuyTransaction) Deserialize(reader *stream.StreamReader) {
	t.Unknown = reader.Read16()
	t.ItemID = reader.ReadU32()
	t.Quantity = reader.ReadU16()

}

type SellTransaction struct {
	Slot     int16
	ItemID   uint32
	Quantity uint16
}

func (t *SellTransaction) GetMode() byte {
	return 1
}

func (t *SellTransaction) Deserialize(reader *stream.StreamReader) {
	t.Slot = reader.Read16()
	t.ItemID = reader.ReadU32()
	t.Quantity = reader.ReadU16()

}

type RechargeTransaction struct {
	Slot int16
}

func (t *RechargeTransaction) GetMode() byte {
	return 2
}

func (t *RechargeTransaction) Deserialize(reader *stream.StreamReader) {
	t.Slot = reader.Read16()

}

func (*NpcShop) Opcode() byte { return 0x2C }

func (n *NpcShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (n *NpcShop) Deserialize(reader *stream.StreamReader) {
	n.Mode = reader.ReadU8()

	var transaction ShopTransaction
	switch n.Mode {
	case 0:
		transaction = &BuyTransaction{}
	case 1:
		transaction = &SellTransaction{}
	case 2:
		transaction = &RechargeTransaction{}
	default:
		return
	}

	transaction.Deserialize(reader)
	n.Transaction = transaction
}
