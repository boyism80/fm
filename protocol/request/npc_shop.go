package request

import "github.com/boyism80/fm/stream"

type ShopTransaction interface {
	GetMode() byte
	Deserialize(reader *stream.StreamReader) error
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

func (t *BuyTransaction) Deserialize(reader *stream.StreamReader) error {
	unknown, err := reader.Read16()
	if err != nil {
		return err
	}
	t.Unknown = unknown

	itemID, err := reader.ReadU32()
	if err != nil {
		return err
	}
	t.ItemID = itemID

	quantity, err := reader.ReadU16()
	if err != nil {
		return err
	}
	t.Quantity = quantity

	return nil
}

type SellTransaction struct {
	Slot     int16
	ItemID   uint32
	Quantity uint16
}

func (t *SellTransaction) GetMode() byte {
	return 1
}

func (t *SellTransaction) Deserialize(reader *stream.StreamReader) error {
	slot, err := reader.Read16()
	if err != nil {
		return err
	}
	t.Slot = slot

	itemID, err := reader.ReadU32()
	if err != nil {
		return err
	}
	t.ItemID = itemID

	quantity, err := reader.ReadU16()
	if err != nil {
		return err
	}
	t.Quantity = quantity

	return nil
}

type RechargeTransaction struct {
	Slot int16
}

func (t *RechargeTransaction) GetMode() byte {
	return 2
}

func (t *RechargeTransaction) Deserialize(reader *stream.StreamReader) error {
	slot, err := reader.Read16()
	if err != nil {
		return err
	}
	t.Slot = slot

	return nil
}

func (n *NpcShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (n *NpcShop) Deserialize(reader *stream.StreamReader) error {
	mode, err := reader.ReadU8()
	if err != nil {
		return err
	}
	n.Mode = mode

	var transaction ShopTransaction
	switch mode {
	case 0:
		transaction = &BuyTransaction{}
	case 1:
		transaction = &SellTransaction{}
	case 2:
		transaction = &RechargeTransaction{}
	default:
		return nil
	}

	if err := transaction.Deserialize(reader); err != nil {
		return err
	}

	n.Transaction = transaction
	return nil
}
