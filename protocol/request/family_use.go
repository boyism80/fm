package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyUse struct {
	Feature internal.FamilyFeature
	Name    string
}

func (*FamilyUse) Opcode() byte { return 0x80 }

func (p *FamilyUse) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(uint32(p.Feature))
	if p.Feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT || p.Feature == internal.FamilyFeature_FAMILY_FEATURE_SUMMON {
		w.WriteStr16(p.Name)
	}
	return nil
}

func (p *FamilyUse) Deserialize(r *stream.StreamReader) {
	p.Feature = internal.FamilyFeature(r.ReadU32())
	if p.Feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT || p.Feature == internal.FamilyFeature_FAMILY_FEATURE_SUMMON {
		p.Name = r.ReadStr16()
	}
}
