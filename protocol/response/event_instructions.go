package response

import "github.com/boyism80/fm/stream"

type EventInstructions struct{}

func (p *EventInstructions) Opcode() uint16 {
	return 0x64
}

func (p *EventInstructions) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(0)
	return nil
}

func (p *EventInstructions) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
}
