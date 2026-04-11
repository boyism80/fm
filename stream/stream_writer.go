package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/util"
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

func (sw *StreamWriter) Write(data []byte) {
	_, err := sw.buf.Write(data)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) WriteU8(val uint8) {
	if err := sw.buf.WriteByte(val); err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) Write8(val int8) {
	sw.WriteU8(uint8(val))
}

func (sw *StreamWriter) WriteU16(val uint16) {
	b := make([]byte, 2)
	sw.order.PutUint16(b, val)
	_, err := sw.buf.Write(b)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) Write16(val int16) {
	sw.WriteU16(uint16(val))
}

func (sw *StreamWriter) WriteU32(val uint32) {
	b := make([]byte, 4)
	sw.order.PutUint32(b, val)
	_, err := sw.buf.Write(b)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) Write32(val int32) {
	sw.WriteU32(uint32(val))
}

func (sw *StreamWriter) WriteU64(val uint64) {
	b := make([]byte, 8)
	sw.order.PutUint64(b, val)
	_, err := sw.buf.Write(b)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) Write64(val int64) {
	sw.WriteU64(uint64(val))
}

func (sw *StreamWriter) WriteFloat32(val float32) {
	sw.WriteU32(math.Float32bits(val))
}

func (sw *StreamWriter) WriteFloat64(val float64) {
	sw.WriteU64(math.Float64bits(val))
}

func (sw *StreamWriter) WriteStr8(str string) {
	if len(str) > 255 {
		panic(errors.New("string too long for WriteStr8"))
	}
	sw.WriteU8(uint8(len(str)))
	_, err := sw.buf.Write([]byte(str))
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) WriteStr16(str string) {
	if len(str) > 65535 {
		panic(errors.New("string too long for WriteStr16"))
	}

	encoder := korean.EUCKR.NewEncoder()
	encodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		panic(err)
	}

	sw.WriteU16(uint16(len(encodedStr)))

	_, err = sw.buf.Write(encodedStr)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) WriteStr32(str string) {
	if len(str) > int(^uint32(0)) {
		panic(errors.New("string too long for WriteStr32"))
	}

	encoder := korean.EUCKR.NewEncoder()
	decodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		panic(err)
	}

	sw.WriteU32(uint32(len(decodedStr)))

	_, err = sw.buf.Write(decodedStr)
	if err != nil {
		panic(err)
	}
}

func (sw *StreamWriter) WriteStaticStr(str string, max int) {
	if len(str) > int(^uint32(0)) {
		panic(errors.New("string too long for WriteStr32"))
	}

	encoder := korean.EUCKR.NewEncoder()
	encodedStr, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		panic(err)
	}

	if len(encodedStr) > max {
		encodedStr = encodedStr[:max]
	}

	sw.Write(encodedStr)
	sw.Write(make([]byte, max-len(encodedStr)))
}

func (sw *StreamWriter) WriteBoolean(val bool) {
	var byteVal uint8
	if val {
		byteVal = 1
	} else {
		byteVal = 0
	}
	sw.WriteU8(byteVal)
}

func (sw *StreamWriter) WriteIPAddress(ip string) {
	if ip == "localhost" {
		ip = "127.0.0.1"
	}
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		panic(errors.New("invalid IP address format"))
	}

	serverIP := make([]byte, 4)
	for i, part := range parts {
		num, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			panic(err)
		}
		serverIP[i] = byte(num)
	}
	sw.Write(serverIP)
}

func (sw *StreamWriter) WriteDateTime(dt time.Time) {
	sw.WriteU64(util.ToFileTime(dt))
}

func (sw *StreamWriter) ToString() string {
	return util.ToHexString(sw.Bytes())
}
