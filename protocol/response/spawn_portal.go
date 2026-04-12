package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

const DisabledPortalMapID uint32 = 999999999

type SpawnPortal struct {
	DestMapID   uint32
	SourceMapID uint32
	SkillID     uint32
	Position    *types.Vector2[int16]
}

func (p *SpawnPortal) Opcode() uint16 {
	return 0x32
}

func (p *SpawnPortal) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.DestMapID)
	w.WriteU32(p.SourceMapID)
	if p.DestMapID != DisabledPortalMapID && p.SourceMapID != DisabledPortalMapID {
		w.WriteU32(p.SkillID)
		if p.Position == nil {
			w.Write16(0)
			w.Write16(0)
			return nil
		}
		w.Write16(p.Position.X)
		w.Write16(p.Position.Y)
	}
	return nil
}

func (p *SpawnPortal) Deserialize(*stream.StreamReader) {}
