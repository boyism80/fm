package resp

import (
	"fmt"

	"github.com/boyism80/fm/common/stream"
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

func (a *Welcome) Serialize(writer *stream.StreamWriter) error {
	err := writer.WriteU16(uint16(13 + len(version)))
	if err != nil {
		return err
	}

	err = writer.WriteU16(MAGIC)
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
