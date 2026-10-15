package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

// Skill selects the client's macro skill window and its administrator messages instead of the item ones.
const lieDetectorSkillDetail = 2

type LieDetectorQuestion struct {
	Image []byte
	Skill bool
}

func (p *LieDetectorQuestion) Opcode() uint16 {
	return 0x1F
}

func (p *LieDetectorQuestion) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.LieDetectorResultQuestion))
	if p.Skill {
		writer.WriteU8(lieDetectorSkillDetail)
	} else {
		writer.WriteU8(4)
	}
	writer.WriteU8(1)
	writer.WriteU32(uint32(len(p.Image)))
	writer.Write(p.Image)
	return nil
}

func (p *LieDetectorQuestion) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Skill = reader.ReadU8() == lieDetectorSkillDetail
	reader.ReadU8()
	p.Image = reader.Read(int(reader.ReadU32()))
}

type LieDetectorStarted struct {
	Target string
}

func (p *LieDetectorStarted) Opcode() uint16 {
	return 0x1F
}

func (p *LieDetectorStarted) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.LieDetectorResultStarted))
	writer.WriteU8(1)
	writer.WriteStr16(p.Target)
	return nil
}

func (p *LieDetectorStarted) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	p.Target = reader.ReadStr16()
}

type LieDetectorResult struct {
	Result pconst.LieDetectorResult
	Skill  bool
}

func (p *LieDetectorResult) Opcode() uint16 {
	return 0x1F
}

func (p *LieDetectorResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	switch {
	case p.Skill && (p.Result == pconst.LieDetectorResultCaught || p.Result == pconst.LieDetectorResultPassed):
		writer.WriteU8(lieDetectorSkillDetail)
	case p.Result == pconst.LieDetectorResultCaught:
		writer.WriteU8(4)
	case p.Result == pconst.LieDetectorResultPassed:
		writer.WriteU8(1)
	default:
		writer.WriteU8(0)
	}
	return nil
}

func (p *LieDetectorResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.LieDetectorResult(reader.ReadU8())
	p.Skill = reader.ReadU8() == lieDetectorSkillDetail && (p.Result == pconst.LieDetectorResultCaught || p.Result == pconst.LieDetectorResultPassed)
}
