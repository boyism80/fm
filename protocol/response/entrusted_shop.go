package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type EntrustedShopCheckResult struct {
	Result  pconst.EntrustedShopCheck
	MapID   uint32
	Channel uint8
}

func (p *EntrustedShopCheckResult) Opcode() uint16 {
	return 0x26
}

func (p *EntrustedShopCheckResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	if p.Result == pconst.EntrustedShopAlreadyOpen {
		writer.WriteU32(p.MapID)
		writer.WriteU8(p.Channel)
	}
	return nil
}

func (p *EntrustedShopCheckResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.EntrustedShopCheck(reader.ReadU8())
	if p.Result == pconst.EntrustedShopAlreadyOpen {
		p.MapID = reader.ReadU32()
		p.Channel = reader.ReadU8()
	}
}

type EntrustedShopBalloon struct {
	SN     uint32
	Title  string
	ItemID uint32
	Users  uint8
}

func (b *EntrustedShopBalloon) serialize(writer *stream.StreamWriter) {
	writer.WriteU8(pconst.MiniRoomTypeEntrustedShop)
	writer.WriteU32(b.SN)
	writer.WriteStr16(b.Title)
	writer.WriteU8(uint8(b.ItemID % 10))
	writer.WriteU8(b.Users)
	writer.WriteU8(pconst.MiniRoomEntrustedShopUsers)
}

func (b *EntrustedShopBalloon) deserialize(reader *stream.StreamReader) {
	if reader.ReadU8() == 0 {
		return
	}
	b.SN = reader.ReadU32()
	b.Title = reader.ReadStr16()
	reader.ReadU8()
	b.Users = reader.ReadU8()
	reader.ReadU8()
}

type SpawnEntrustedShop struct {
	EmployerID uint32
	X          int16
	Y          int16
	Foothold   uint16
	OwnerName  string
	EntrustedShopBalloon
}

func (p *SpawnEntrustedShop) Opcode() uint16 {
	return 0xC3
}

func (p *SpawnEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.EmployerID)
	writer.WriteU32(p.ItemID)
	writer.Write16(p.X)
	writer.Write16(p.Y)
	writer.WriteU16(p.Foothold)
	writer.WriteStr16(p.OwnerName)
	p.serialize(writer)
	return nil
}

func (p *SpawnEntrustedShop) Deserialize(reader *stream.StreamReader) {
	p.EmployerID = reader.ReadU32()
	p.ItemID = reader.ReadU32()
	p.X = reader.Read16()
	p.Y = reader.Read16()
	p.Foothold = reader.ReadU16()
	p.OwnerName = reader.ReadStr16()
	p.deserialize(reader)
}

type DestroyEntrustedShop struct {
	EmployerID uint32
}

func (p *DestroyEntrustedShop) Opcode() uint16 {
	return 0xC4
}

func (p *DestroyEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.EmployerID)
	return nil
}

func (p *DestroyEntrustedShop) Deserialize(reader *stream.StreamReader) {
	p.EmployerID = reader.ReadU32()
}

type UpdateEntrustedShop struct {
	EmployerID uint32
	EntrustedShopBalloon
}

func (p *UpdateEntrustedShop) Opcode() uint16 {
	return 0xC5
}

func (p *UpdateEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.EmployerID)
	p.serialize(writer)
	return nil
}

func (p *UpdateEntrustedShop) Deserialize(reader *stream.StreamReader) {
	p.EmployerID = reader.ReadU32()
	p.deserialize(reader)
}

type StoreBankOpen struct {
	NpcID uint32
	Storage
}

func (p *StoreBankOpen) Opcode() uint16 {
	return 0xEC
}

func (p *StoreBankOpen) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.StoreBankResultOpen))
	writer.WriteU32(p.NpcID)
	p.serializeTabs(writer)
	return nil
}

func (p *StoreBankOpen) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NpcID = reader.ReadU32()
	p.deserializeTabs(reader)
}

type StoreBankFee struct {
	Days uint32
	Fee  int32
}

func (p *StoreBankFee) Opcode() uint16 {
	return 0xEC
}

func (p *StoreBankFee) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.StoreBankResultFee))
	writer.WriteU32(p.Days)
	writer.Write32(p.Fee)
	return nil
}

func (p *StoreBankFee) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Days = reader.ReadU32()
	p.Fee = reader.Read32()
}

type StoreBankLocation struct {
	NpcID   uint32
	MapID   uint32
	Channel uint8
}

func (p *StoreBankLocation) Opcode() uint16 {
	return 0xEC
}

func (p *StoreBankLocation) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.StoreBankResultLocation))
	writer.WriteU32(p.NpcID)
	writer.WriteU32(p.MapID)
	writer.WriteU8(p.Channel)
	return nil
}

func (p *StoreBankLocation) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NpcID = reader.ReadU32()
	p.MapID = reader.ReadU32()
	p.Channel = reader.ReadU8()
}

type StoreBankResult struct {
	Result pconst.StoreBankResult
}

func (p *StoreBankResult) Opcode() uint16 {
	return 0xEB
}

func (p *StoreBankResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	return nil
}

func (p *StoreBankResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.StoreBankResult(reader.ReadU8())
}
