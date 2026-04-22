package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfItemMakerSuccessEffect struct {
}

func (p *ShowSelfItemMakerSuccessEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfItemMakerSuccessEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0x11)
	writer.WriteU32(0)
	return nil
}

func (p *ShowSelfItemMakerSuccessEffect) Deserialize(reader *stream.StreamReader) {
}
