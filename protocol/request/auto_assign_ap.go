package request

import (
	"github.com/boyism80/fm/stream"
)

// AutoAssignAP represents the AUTO_ASSIGN_AP packet (0x47)
// Used to distribute multiple AP to two stats simultaneously (STR, DEX, INT, LUK only)
type AutoAssignAP struct {
	Tick          uint32
	Unknown       uint32
	PrimaryStat   uint32 // 64=STR, 128=DEX, 256=INT, 512=LUK
	Amount        uint32 // Amount of AP to invest in primary stat
	SecondaryStat uint32 // 64=STR, 128=DEX, 256=INT, 512=LUK
	Amount2       uint32 // Amount of AP to invest in secondary stat
}

func (p *AutoAssignAP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *AutoAssignAP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Unknown = reader.ReadU32()

	if reader.Remaining() < 16 {
		return
	}

	p.PrimaryStat = reader.ReadU32()
	p.Amount = reader.ReadU32()
	p.SecondaryStat = reader.ReadU32()
	p.Amount2 = reader.ReadU32()
}
