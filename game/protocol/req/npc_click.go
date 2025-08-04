package req

import "github.com/boyism80/fm/core/stream"

type NpcClick struct {
	OID uint32
}

func (n *NpcClick) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (n *NpcClick) Deserialize(reader *stream.StreamReader) error {
	oid, err := reader.ReadU32()
	if err != nil {
		return err
	}
	n.OID = oid
	return nil
}
