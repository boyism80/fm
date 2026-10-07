package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type CashShopItem struct {
	Serial      uint64
	AccountID   uint32
	CharacterID uint32
	ItemID      uint32
	CommoditySN uint32
	Count       uint16
	BuyerName   string
	Expiration  time.Time
}

func (i *CashShopItem) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(i.Serial)
	writer.WriteU32(i.AccountID)
	writer.WriteU32(i.CharacterID)
	writer.WriteU32(i.ItemID)
	writer.WriteU32(i.CommoditySN)
	writer.WriteU16(i.Count)
	writer.WriteStaticStr(i.BuyerName, 13)
	writer.WriteDateTime(i.Expiration)
	writer.WriteU32(0)
	writer.WriteU32(0)
}

func (i *CashShopItem) Deserialize(reader *stream.StreamReader) {
	i.Serial = reader.ReadU64()
	i.AccountID = reader.ReadU32()
	i.CharacterID = reader.ReadU32()
	i.ItemID = reader.ReadU32()
	i.CommoditySN = reader.ReadU32()
	i.Count = reader.ReadU16()
	i.BuyerName = reader.ReadStaticStr(13)
	i.Expiration = util.FromFileTime(reader.ReadU64())
	reader.Skip(8)
}

type CashShopGift struct {
	Serial     uint64
	ItemID     uint32
	SenderName string
	Message    string
}

func (g *CashShopGift) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(g.Serial)
	writer.WriteU32(g.ItemID)
	writer.WriteStaticStr(g.SenderName, 13)
	writer.WriteStaticStr(g.Message, 73)
}

func (g *CashShopGift) Deserialize(reader *stream.StreamReader) {
	g.Serial = reader.ReadU64()
	g.ItemID = reader.ReadU32()
	g.SenderName = reader.ReadStaticStr(13)
	g.Message = reader.ReadStaticStr(73)
}

type CashShopGrantedItem struct {
	Count  uint16
	Slot   uint16
	ItemID uint32
}

func (g *CashShopGrantedItem) Serialize(writer *stream.StreamWriter) {
	writer.WriteU16(g.Count)
	writer.WriteU16(g.Slot)
	writer.WriteU32(g.ItemID)
}

func (g *CashShopGrantedItem) Deserialize(reader *stream.StreamReader) {
	g.Count = reader.ReadU16()
	g.Slot = reader.ReadU16()
	g.ItemID = reader.ReadU32()
}
