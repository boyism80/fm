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

func (p *AutoAssignAP) Deserialize(reader *stream.StreamReader) error {
	var err error
	if p.Tick, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.Unknown, err = reader.ReadU32(); err != nil {
		return err
	}

	// Check if we have enough data (16 bytes for PrimaryStat, Amount, SecondaryStat, Amount2)
	if reader.Remaining() < 16 {
		return nil // Not enough data, but don't error
	}

	if p.PrimaryStat, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.Amount, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.SecondaryStat, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.Amount2, err = reader.ReadU32(); err != nil {
		return err
	}
	return nil
}
