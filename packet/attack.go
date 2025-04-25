package packet

import "github.com/boyism80/fm/stream"

// Attack 패킷 타입
type AttackPacket struct {
	AttackerID int32 // 공격자 ID
	TargetID   int32 // 대상 ID
}

// Serialize는 AttackPacket을 StreamWriter를 사용해 직렬화
func (a *AttackPacket) Serialize(writer *stream.StreamWriter) error {
	err := writer.Write32(a.AttackerID)
	if err != nil {
		return err
	}
	err = writer.Write32(a.TargetID)
	return err
}

// Deserialize는 StreamReader를 사용해 AttackPacket을 역직렬화
func (a *AttackPacket) Deserialize(reader *stream.StreamReader) error {
	attackerID, err := reader.Read32()
	if err != nil {
		return err
	}
	targetID, err := reader.Read32()
	if err != nil {
		return err
	}
	a.AttackerID = attackerID
	a.TargetID = targetID
	return nil
}
