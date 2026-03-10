package request

import (
	"github.com/boyism80/fm/stream"
)

type ChangeKeymap struct {
	// Mode 1: Keymap changes (available > 8)
	Skip       int32
	NumChanges int32
	Changes    []KeymapChange

	// Mode 2: Pet auto pot (available <= 8)
	Type int32
	Data int32
}

type KeymapChange struct {
	Key    int32
	Type   byte
	Action int32
}

func (c *ChangeKeymap) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *ChangeKeymap) Deserialize(reader *stream.StreamReader) error {
	available := reader.Remaining()
	if available > 8 {
		// Mode 1: Keymap changes
		var err error
		if c.Skip, err = reader.Read32(); err != nil {
			return err
		}
		if c.NumChanges, err = reader.Read32(); err != nil {
			return err
		}
		c.Changes = make([]KeymapChange, c.NumChanges)
		for i := int32(0); i < c.NumChanges; i++ {
			if c.Changes[i].Key, err = reader.Read32(); err != nil {
				return err
			}
			if c.Changes[i].Type, err = reader.ReadU8(); err != nil {
				return err
			}
			if c.Changes[i].Action, err = reader.Read32(); err != nil {
				return err
			}
		}
	} else {
		// Mode 2: Pet auto pot
		var err error
		if c.Type, err = reader.Read32(); err != nil {
			return err
		}
		if c.Data, err = reader.Read32(); err != nil {
			return err
		}
	}
	return nil
}
