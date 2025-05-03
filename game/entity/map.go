package entity

import "github.com/boyism80/fm/common/types"

type Map struct {
	Id          uint32                      // 맵 ID
	Name        string                      // 맵 이름
	SpawnPoints map[uint8]types.Vec2[int32] // 스폰 위치
	Bounds      types.Rect[int32]           // 맵 범위 (좌표계)
	IsTown      bool                        // 마을 여부
	HasClock    bool                        // 시계 UI
	Properties  map[string]string           // 기타 설정 (BGM, weather 등)
}
