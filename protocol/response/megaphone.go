package response

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Megaphone struct {
	Type    constant.ServerMessageType
	Lines   []string
	Channel int
	Ear     bool
	Item    dto.Item
}

func (*Megaphone) Opcode() uint16 { return 0x33 }

func (p *Megaphone) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Type))
	w.WriteStr16(p.Lines[0])
	if p.Type == constant.MsgThreeMegaphoneLines {
		w.WriteU8(uint8(len(p.Lines)))
		for _, line := range p.Lines[1:] {
			w.WriteStr16(line)
		}
	}
	w.WriteU8(uint8(max(p.Channel, 1) - 1))
	w.WriteBoolean(p.Ear)
	if p.Type == constant.MsgItemMegaphone {
		w.WriteBoolean(p.Item != nil)
		if p.Item != nil {
			p.Item.Serialize(w, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
		}
	}
	return nil
}

func (p *Megaphone) Deserialize(r *stream.StreamReader) {
	p.Type = constant.ServerMessageType(r.ReadU8())
	p.Lines = []string{r.ReadStr16()}
	if p.Type == constant.MsgThreeMegaphoneLines {
		count := int(r.ReadU8())
		for i := 1; i < count; i++ {
			p.Lines = append(p.Lines, r.ReadStr16())
		}
	}
	p.Channel = int(r.ReadU8()) + 1
	p.Ear = r.ReadBool()
	if p.Type == constant.MsgItemMegaphone && r.ReadBool() {
		p.Item = dto.NewItemFromStream(r)
	}
}
