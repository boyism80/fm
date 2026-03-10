package request

import (
	"github.com/boyism80/fm/stream"
)

// DistributeSP represents the DISTRIBUTE_SP packet (0x49)
// Used to distribute 1 SP to a skill
type DistributeSP struct {
	Tick    uint32
	SkillID uint32
}

func (p *DistributeSP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DistributeSP) Deserialize(reader *stream.StreamReader) error {
	var err error
	if p.Tick, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.SkillID, err = reader.ReadU32(); err != nil {
		return err
	}
	return nil
}
