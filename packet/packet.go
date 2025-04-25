package packet

import "github.com/boyism80/fm/stream"

// Packet 인터페이스 정의
type Packet interface {
	Serialize(writer *stream.StreamWriter) error   // 패킷을 직렬화
	Deserialize(reader *stream.StreamReader) error // 패킷을 역직렬화
}
