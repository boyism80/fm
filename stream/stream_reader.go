package stream

import (
	"encoding/binary"
	"errors"
	"math"
)

type Endian binary.ByteOrder

var (
	LittleEndian Endian = Endian(binary.LittleEndian)
	BigEndian    Endian = Endian(binary.BigEndian)
)

type StreamReader struct {
	data   []byte
	cursor int
	order  binary.ByteOrder
}

func NewStreamReader(data []byte, endian Endian) *StreamReader {
	return &StreamReader{
		data:   data,
		cursor: 0,
		order:  binary.ByteOrder(endian),
	}
}

func (sr *StreamReader) Read(n int) ([]byte, error) {
	if sr.cursor+n > len(sr.data) {
		return nil, errors.New("not enough bytes to read")
	}
	b := sr.data[sr.cursor : sr.cursor+n]
	sr.cursor += n
	return b, nil
}

func (sr *StreamReader) ReadU8() (uint8, error) {
	b, err := sr.Read(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (sr *StreamReader) Read8() (int8, error) {
	val, err := sr.ReadU8()
	return int8(val), err
}

func (sr *StreamReader) ReadU16() (uint16, error) {
	b, err := sr.Read(2)
	if err != nil {
		return 0, err
	}
	return sr.order.Uint16(b), nil
}

func (sr *StreamReader) Read16() (int16, error) {
	val, err := sr.ReadU16()
	return int16(val), err
}

func (sr *StreamReader) ReadU32() (uint32, error) {
	b, err := sr.Read(4)
	if err != nil {
		return 0, err
	}
	return sr.order.Uint32(b), nil
}

func (sr *StreamReader) Read32() (int32, error) {
	val, err := sr.ReadU32()
	return int32(val), err
}

func (sr *StreamReader) ReadU64() (uint64, error) {
	b, err := sr.Read(8)
	if err != nil {
		return 0, err
	}
	return sr.order.Uint64(b), nil
}

func (sr *StreamReader) Read64() (int64, error) {
	val, err := sr.ReadU64()
	return int64(val), err
}

func (sr *StreamReader) ReadFloat32() (float32, error) {
	bits, err := sr.ReadU32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(bits), nil
}

func (sr *StreamReader) ReadFloat64() (float64, error) {
	bits, err := sr.ReadU64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bits), nil
}

func (sr *StreamReader) ReadStr8() (string, error) {
	length, err := sr.ReadU8()
	if err != nil {
		return "", err
	}
	b, err := sr.Read(int(length))
	return string(b), err
}

func (sr *StreamReader) ReadStr16() (string, error) {
	length, err := sr.ReadU16()
	if err != nil {
		return "", err
	}
	b, err := sr.Read(int(length))
	return string(b), err
}

func (sr *StreamReader) ReadStr32() (string, error) {
	length, err := sr.ReadU32()
	if err != nil {
		return "", err
	}
	b, err := sr.Read(int(length))
	return string(b), err
}

func (sr *StreamReader) Reset() {
	sr.cursor = 0
}

func (sr *StreamReader) DiscardRead() {
	if sr.cursor > 0 {
		sr.data = sr.data[sr.cursor:]
		sr.cursor = 0
	}
}

func (sr *StreamReader) Remaining() int {
	return len(sr.data) - sr.cursor
}
