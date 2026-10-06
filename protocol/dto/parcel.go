package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type Parcel struct {
	ID      uint32
	Sender  string
	Meso    int32
	Expire  time.Time
	Quick   bool
	Message string
	Item    Item
}

func (p *Parcel) Serialize(writer *stream.StreamWriter) {
	writer.WriteU32(p.ID)
	writer.WriteStaticStr(p.Sender, 13)
	writer.Write32(p.Meso)
	writer.WriteDateTime(p.Expire)
	if p.Quick {
		writer.Write32(1)
	} else {
		writer.Write32(0)
	}
	writer.WriteStaticStr(p.Message, 200)
	writer.WriteU8(0)
	if p.Item == nil {
		writer.WriteU8(0)
		return
	}
	writer.WriteU8(1)
	p.Item.Serialize(writer, ItemSerializeOption{SlotMode: SlotEncodeOmit})
}

func NewParcelFromStream(reader *stream.StreamReader) Parcel {
	p := Parcel{}
	p.ID = reader.ReadU32()
	p.Sender = reader.ReadStaticStr(13)
	p.Meso = reader.Read32()
	p.Expire = util.FromFileTime(reader.ReadU64())
	p.Quick = reader.Read32() != 0
	p.Message = reader.ReadStaticStr(200)
	reader.ReadU8()
	if reader.ReadU8() != 0 {
		p.Item = NewItemFromStream(reader)
	}
	return p
}
