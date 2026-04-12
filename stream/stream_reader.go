package stream

import (
	"encoding/binary"
	"errors"
	"math"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

var ErrTruncated = errors.New("stream: not enough bytes")

var ErrSkipOutOfBounds = errors.New("stream: skip out of bounds")

type Endian binary.ByteOrder

var (
	LittleEndian Endian = Endian(binary.LittleEndian)
	BigEndian    Endian = Endian(binary.BigEndian)
)

type StreamReader struct {
	data   *[]byte
	cursor int
	order  binary.ByteOrder
}

func NewStreamReader(data *[]byte, endian Endian) *StreamReader {
	return &StreamReader{
		data:   data,
		cursor: 0,
		order:  binary.ByteOrder(endian),
	}
}

func (sr *StreamReader) Read(n int) []byte {
	if n < 0 {
		panic("stream: negative read length")
	}
	if sr.cursor+n > len(*sr.data) {
		panic(ErrTruncated)
	}
	bytes := *sr.data
	b := bytes[sr.cursor : sr.cursor+n]
	sr.cursor += n
	return b
}

func (sr *StreamReader) ReadBool() bool {
	b := sr.ReadU8()
	return b != 0
}

func (sr *StreamReader) ReadU8() uint8 {
	b := sr.Read(1)
	return b[0]
}

func (sr *StreamReader) Read8() int8 {
	return int8(sr.ReadU8())
}

func (sr *StreamReader) ReadU16() uint16 {
	b := sr.Read(2)
	return sr.order.Uint16(b)
}

func (sr *StreamReader) Read16() int16 {
	return int16(sr.ReadU16())
}

func (sr *StreamReader) ReadU32() uint32 {
	b := sr.Read(4)
	return sr.order.Uint32(b)
}

func (sr *StreamReader) Read32() int32 {
	return int32(sr.ReadU32())
}

func (sr *StreamReader) ReadU64() uint64 {
	b := sr.Read(8)
	return sr.order.Uint64(b)
}

func (sr *StreamReader) Read64() int64 {
	return int64(sr.ReadU64())
}

func (sr *StreamReader) ReadFloat32() float32 {
	return math.Float32frombits(sr.ReadU32())
}

func (sr *StreamReader) ReadFloat64() float64 {
	return math.Float64frombits(sr.ReadU64())
}

func (sr *StreamReader) ReadStr8() string {
	length := int(sr.ReadU8())
	return string(sr.Read(length))
}

func (sr *StreamReader) ReadStr16() string {
	length := int(sr.ReadU16())
	b := sr.Read(length)
	decoder := korean.EUCKR.NewDecoder()
	decodedStr, _, err := transform.Bytes(decoder, b)
	if err != nil {
		panic(err)
	}
	return string(decodedStr)
}

func (sr *StreamReader) ReadStr32() string {
	length := int(sr.ReadU32())
	b := sr.Read(length)
	decoder := korean.EUCKR.NewDecoder()
	decodedStr, _, err := transform.Bytes(decoder, b)
	if err != nil {
		panic(err)
	}
	return string(decodedStr)
}

func (sr *StreamReader) Reset() {
	sr.cursor = 0
}

func (sr *StreamReader) DiscardRead() {
	if sr.cursor > 0 {
		bytes := *sr.data
		*sr.data = bytes[sr.cursor:]
		sr.cursor = 0
	}
}

func (sr *StreamReader) Remaining() int {
	return len(*sr.data) - sr.cursor
}

func (sr *StreamReader) Skip(n int) {
	if n < 0 {
		panic("stream: negative skip length")
	}
	if n > sr.Remaining() {
		panic(ErrSkipOutOfBounds)
	}
	sr.cursor += n
}
