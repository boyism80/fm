package response

import "github.com/boyism80/fm/stream"

const (
	SecondPasswordResultSuccess = 0
	SecondPasswordResultInvalid = 0x14
)

type SecondPasswordResult struct {
	Mode   uint8
	Result uint8
}

func (*SecondPasswordResult) Opcode() uint16 { return 0x0F }

func (a *SecondPasswordResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(a.Mode)
	writer.WriteU8(a.Result)
	return nil
}

func (a *SecondPasswordResult) Deserialize(reader *stream.StreamReader) {
	a.Mode = reader.ReadU8()
	a.Result = reader.ReadU8()
}
