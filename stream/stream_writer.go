package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

type StreamWriter struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

func NewStreamWriter(endian Endian) *StreamWriter {
	return &StreamWriter{
		buf:   new(bytes.Buffer),
		order: binary.ByteOrder(endian),
	}
}

func (sw *StreamWriter) Bytes() []byte {
	return sw.buf.Bytes()
}

func (sw *StreamWriter) Write(data []byte) error {
	_, err := sw.buf.Write(data)
	return err
}

func (sw *StreamWriter) WriteU8(val uint8) error {
	return sw.buf.WriteByte(val)
}

func (sw *StreamWriter) Write8(val int8) error {
	return sw.WriteU8(uint8(val))
}

func (sw *StreamWriter) WriteU16(val uint16) error {
	b := make([]byte, 2)
	sw.order.PutUint16(b, val)
	_, err := sw.buf.Write(b)
	return err
}

func (sw *StreamWriter) Write16(val int16) error {
	return sw.WriteU16(uint16(val))
}

func (sw *StreamWriter) WriteU32(val uint32) error {
	b := make([]byte, 4)
	sw.order.PutUint32(b, val)
	_, err := sw.buf.Write(b)
	return err
}

func (sw *StreamWriter) Write32(val int32) error {
	return sw.WriteU32(uint32(val))
}

func (sw *StreamWriter) WriteU64(val uint64) error {
	b := make([]byte, 8)
	sw.order.PutUint64(b, val)
	_, err := sw.buf.Write(b)
	return err
}

func (sw *StreamWriter) Write64(val int64) error {
	return sw.WriteU64(uint64(val))
}

func (sw *StreamWriter) WriteFloat32(val float32) error {
	return sw.WriteU32(math.Float32bits(val))
}

func (sw *StreamWriter) WriteFloat64(val float64) error {
	return sw.WriteU64(math.Float64bits(val))
}

func (sw *StreamWriter) WriteStr8(str string) error {
	if len(str) > 255 {
		return errors.New("string too long for WriteStr8")
	}
	if err := sw.WriteU8(uint8(len(str))); err != nil {
		return err
	}
	_, err := sw.buf.Write([]byte(str))
	return err
}

func (sw *StreamWriter) WriteStr16(str string) error {
	if len(str) > 65535 {
		return errors.New("string too long for WriteStr16")
	}

	encoder := korean.EUCKR.NewEncoder()
	encodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return err
	}

	if err := sw.WriteU16(uint16(len(encodedStr))); err != nil {
		return err
	}

	_, err = sw.buf.Write(encodedStr)
	return err
}

func (sw *StreamWriter) WriteStr32(str string) error {
	if len(str) > int(^uint32(0)) {
		return errors.New("string too long for WriteStr32")
	}

	encoder := korean.EUCKR.NewEncoder()
	decodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return err
	}

	if err := sw.WriteU32(uint32(len(decodedStr))); err != nil {
		return err
	}

	_, err = sw.buf.Write(decodedStr)
	return err
}

func (sw *StreamWriter) WriteStaticStr(str string, max int) error {
	if len(str) > int(^uint32(0)) {
		return errors.New("string too long for WriteStr32")
	}

	encoder := korean.EUCKR.NewEncoder()
	encodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return err
	}

	if len(encodedStr) > max {
		encodedStr = encodedStr[:max]
	}

	sw.Write(encodedStr)
	sw.Write(make([]byte, max-len(encodedStr)))
	return nil
}

func (sw *StreamWriter) WriteBoolean(val bool) error {
	var byteVal uint8
	if val {
		byteVal = 1
	} else {
		byteVal = 0
	}
	return sw.WriteU8(byteVal)
}

func (sw *StreamWriter) WriteIPAddress(ip string) error {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return errors.New("invalid IP address format")
	}

	serverIP := make([]byte, 4)
	for i, part := range parts {
		num, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return err
		}
		serverIP[i] = byte(num)
	}
	return sw.Write(serverIP)
}
