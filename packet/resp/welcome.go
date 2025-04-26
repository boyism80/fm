package resp

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

var version string

const (
	MAPLE_VERSION = 10
	MAPLE_CHECK   = 1
	MAPLE_PATCH   = 1
	MAPLE_MAGIC   = 291
)

func init() {
	ret := 0
	ret ^= (MAPLE_VERSION & 0x7FFF)
	ret ^= (MAPLE_CHECK << 15)
	ret ^= ((MAPLE_PATCH & 0xFF) << 16)
	version = fmt.Sprintf("%d", ret)
}

type Welcome struct {
	SendIv []byte
	RecvIv []byte
}

func (a *Welcome) Serialize(writer *stream.StreamWriter) error {
	err := writer.WriteU16(uint16(13 + len(version)))
	if err != nil {
		return err
	}

	err = writer.WriteU16(MAPLE_MAGIC)
	if err != nil {
		return err
	}
	err = writer.WriteStr16(version)
	if err != nil {
		return err
	}
	err = writer.Write(a.RecvIv)
	if err != nil {
		return err
	}
	err = writer.Write(a.SendIv)
	if err != nil {
		return err
	}
	err = writer.WriteU8(1)
	if err != nil {
		return err
	}
	return nil
}

func (a *Welcome) Deserialize(reader *stream.StreamReader) error {
	return nil
}
