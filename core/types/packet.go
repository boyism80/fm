package types

import "github.com/boyism80/fm/core/stream"

type Packet interface {
	Serialize(writer *stream.StreamWriter) error
	Deserialize(reader *stream.StreamReader) error
	Opcode() uint16 // Returns the opcode for this packet type
}
