package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Duey struct {
	Result pconst.DueyResult
}

func (p *Duey) Opcode() uint16 {
	return 0xF7
}

func (p *Duey) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	if p.Result == pconst.DueyResultIdentity {
		writer.WriteU8(0)
	}
	return nil
}

func (p *Duey) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.DueyResult(reader.ReadU8())
	if p.Result == pconst.DueyResultIdentity {
		reader.ReadU8()
	}
}

type DueyOpen struct {
	FromArrival bool
	Parcels     []dto.Parcel
	Expired     []dto.Parcel
}

func (p *DueyOpen) Opcode() uint16 {
	return 0xF7
}

func (p *DueyOpen) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.DueyResultOpen))
	writer.WriteBoolean(p.FromArrival)
	writer.WriteU8(uint8(len(p.Parcels)))
	for i := range p.Parcels {
		p.Parcels[i].Serialize(writer)
	}
	writer.WriteU8(uint8(len(p.Expired)))
	for i := range p.Expired {
		p.Expired[i].Serialize(writer)
	}
	return nil
}

func (p *DueyOpen) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.FromArrival = reader.ReadBool()
	p.Parcels = make([]dto.Parcel, int(reader.ReadU8()))
	for i := range p.Parcels {
		p.Parcels[i] = dto.NewParcelFromStream(reader)
	}
	p.Expired = make([]dto.Parcel, int(reader.ReadU8()))
	for i := range p.Expired {
		p.Expired[i] = dto.NewParcelFromStream(reader)
	}
}

type DueyRemoved struct {
	ParcelID uint32
	Reason   uint8
}

func (p *DueyRemoved) Opcode() uint16 {
	return 0xF7
}

func (p *DueyRemoved) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.DueyResultRemoved))
	writer.WriteU32(p.ParcelID)
	writer.WriteU8(p.Reason)
	return nil
}

func (p *DueyRemoved) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.ParcelID = reader.ReadU32()
	p.Reason = reader.ReadU8()
}

type DueyArrival struct {
	Sender string
	Quick  bool
}

func (p *DueyArrival) Opcode() uint16 {
	return 0xF7
}

func (p *DueyArrival) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.DueyResultArrival))
	writer.WriteStr16(p.Sender)
	writer.WriteBoolean(p.Quick)
	return nil
}

func (p *DueyArrival) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Sender = reader.ReadStr16()
	p.Quick = reader.ReadBool()
}

type DueyArrivals struct {
	Quick bool
}

func (p *DueyArrivals) Opcode() uint16 {
	return 0xF7
}

func (p *DueyArrivals) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.DueyResultArrivals))
	writer.WriteBoolean(p.Quick)
	return nil
}

func (p *DueyArrivals) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Quick = reader.ReadBool()
}
