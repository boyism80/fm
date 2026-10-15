package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type LieDetectorQuestion struct {
	Image []byte
}

func (p *LieDetectorQuestion) Opcode() uint16 {
	return 0x1F
}

func (p *LieDetectorQuestion) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.LieDetectorResultQuestion))
	writer.WriteU8(4)
	writer.WriteU8(1)
	writer.WriteU32(uint32(len(p.Image)))
	writer.Write(p.Image)
	return nil
}

func (p *LieDetectorQuestion) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	reader.ReadU8()
	p.Image = reader.Read(int(reader.ReadU32()))
}

type LieDetectorResult struct {
	Result pconst.LieDetectorResult
}

func (p *LieDetectorResult) Opcode() uint16 {
	return 0x1F
}

func (p *LieDetectorResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	switch p.Result {
	case pconst.LieDetectorResultCaught:
		writer.WriteU8(4)
	case pconst.LieDetectorResultPassed:
		writer.WriteU8(1)
	default:
		writer.WriteU8(0)
	}
	return nil
}

func (p *LieDetectorResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.LieDetectorResult(reader.ReadU8())
	reader.ReadU8()
}
