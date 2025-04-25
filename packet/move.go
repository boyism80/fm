package packet

import "github.com/boyism80/fm/stream"

// Move 패킷 타입
type MovePacket struct {
	X int32 // X 좌표
	Y int32 // Y 좌표
}

// Serialize는 MovePacket을 StreamWriter를 사용해 직렬화
func (m *MovePacket) Serialize(writer *stream.StreamWriter) error {
	err := writer.Write32(m.X)
	if err != nil {
		return err
	}
	err = writer.Write32(m.Y)
	return err
}

// Deserialize는 StreamReader를 사용해 MovePacket을 역직렬화
func (m *MovePacket) Deserialize(reader *stream.StreamReader) error {
	x, err := reader.Read32()
	if err != nil {
		return err
	}
	y, err := reader.Read32()
	if err != nil {
		return err
	}
	m.X = x
	m.Y = y
	return nil
}
