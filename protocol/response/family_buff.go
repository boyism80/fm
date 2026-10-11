package response

import (
	"time"

	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

// ExpRate and DropRate are percentages; 100 means the buff does not raise it.
type FamilyBuff struct {
	Feature   internal.FamilyFeature
	ExpRate   int
	DropRate  int
	Remaining time.Duration
}

func (p *FamilyBuff) Opcode() uint16 {
	return 0x52
}

func (p *FamilyBuff) Serialize(w *stream.StreamWriter) error {
	if p.Remaining <= 0 {
		w.WriteU8(uint8(pconst.FamilyBuffNone))
		return nil
	}

	exp, drop := uint32(0), uint32(0)
	if p.ExpRate > 100 {
		exp = uint32(p.ExpRate)
	}
	if p.DropRate > 100 {
		drop = uint32(p.DropRate)
	}

	kind := pconst.FamilyBuffNone
	switch {
	case exp != 0 && drop != 0:
		kind = pconst.FamilyBuffBoth
	case exp != 0:
		kind = pconst.FamilyBuffExp
	case drop != 0:
		kind = pconst.FamilyBuffDrop
	}
	w.WriteU8(uint8(kind))
	w.WriteU32(uint32(p.Feature))
	w.WriteU32(exp)
	w.WriteU32(drop)
	w.WriteU8(0)
	w.WriteU32(uint32(p.Remaining.Milliseconds()))
	return nil
}

func (p *FamilyBuff) Deserialize(reader *stream.StreamReader) {
}
