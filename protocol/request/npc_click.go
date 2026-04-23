package request

import "github.com/boyism80/fm/stream"

type NpcClick struct {
	OID uint32
}

func (*NpcClick) Opcode() byte { return 0x29 }

func (n *NpcClick) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (n *NpcClick) Deserialize(reader *stream.StreamReader) {
	n.OID = reader.ReadU32()
}
