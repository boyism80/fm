package response

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

var version string

const (
	VERSION = 10
	CHECK   = 1
	PATCH   = 1
	MAGIC   = 291
)

func init() {
	ret := 0
	ret ^= (VERSION & 0x7FFF)
	ret ^= (CHECK << 15)
	ret ^= ((PATCH & 0xFF) << 16)
	version = fmt.Sprintf("%d", ret)
}

type Welcome struct {
	SendIv []byte
	RecvIv []byte
}

// Opcode returns the packet opcode for Welcome
func (a *Welcome) Opcode() uint16 {
	return uint16(13 + len(version))
}

func (a *Welcome) Serialize(writer *stream.StreamWriter) error {
	_ = writer.WriteU16(MAGIC)
	_ = writer.WriteStr16(version)
	_ = writer.Write(a.RecvIv)
	_ = writer.Write(a.SendIv)
	_ = writer.WriteU8(1)
	return nil
}

func (a *Welcome) Deserialize(reader *stream.StreamReader) {
}
