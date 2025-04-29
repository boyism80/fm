package resp

import (
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnPlayer struct {
	Character         *entity.Character
	GuildName         string     // 길드 이름 (없으면 "")
	GuildLogoBG       uint16     // 길드 엠블럼 배경 모양
	GuildLogoBGColor  uint8      // 길드 엠블럼 배경 색
	GuildLogo         uint16     // 길드 엠블럼 모양
	GuildLogoColor    uint8      // 길드 엠블럼 색
	BuffStates        [4]uint32  // 버프 bitmask
	Diseases          [4]uint32  // 디버프 bitmask
	SpeedBuff         uint8      // 버프된 속도 (옵션)
	ComboCount        uint8      // 콤보 카운트 (옵션)
	WKChargeSkillId   uint32     // WK 차지 스킬 ID (옵션)
	MorphId           uint16     // 변신 ID (옵션)
	SpiritClawSkillId uint32     // 스피릿클로 스킬 ID (옵션)
	ItemEffectId      uint32     // 캐릭터 적용 아이템 이펙트
	ChairId           uint32     // 앉아있는 의자 ID
	Balloons          uint32     // 풍선 개수
	Position          types.Vec2 // 캐릭터 실제 위치 (좌표)
	Stance            uint8      // 스탠스
	MountLevel        uint32     // 탈것 레벨
	MountExp          uint32     // 탈것 경험치
	MountFatigue      uint32     // 탈것 피로도
	Chalkboard        string     // 칠판 텍스트 (없으면 "")
	Team              uint8      // 팀 (0/1)
	CrushRings        []*entity.Ring
	FriendshipRings   []*entity.Ring
	MarriageRings     []*entity.Ring
}

func writeRings(writer *stream.StreamWriter, rings []*entity.Ring) error {
	writer.WriteU8(uint8(len(rings))) // 링 개수
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		writer.WriteU64(ring.RingId)       // 고유 링 ID
		writer.WriteU64(ring.PartnerId)    // 상대 캐릭터 ID
		writer.WriteU64(ring.RingUniqueId) // 링 고유 UID (있다면, 없으면 0)
	}
	return nil
}

func (p *SpawnPlayer) Serialize(writer *stream.StreamWriter) error {
	if p.Character == nil {
		return nil
	}

	// 1. 캐릭터 ID
	writer.WriteU32(p.Character.Id)

	// 2. 캐릭터 이름 (MapleAsciiString)
	writer.WriteStr16(p.Character.Name)

	// 3. 길드 정보
	if p.GuildName == "" {
		writer.WriteU32(0)
		writer.WriteU32(0)
	} else {
		writer.WriteStr16(p.GuildName)
		writer.WriteU16(p.GuildLogoBG)
		writer.WriteU8(p.GuildLogoBGColor)
		writer.WriteU16(p.GuildLogo)
		writer.WriteU8(p.GuildLogoColor)
	}

	// 4. Secondary Stats (Buffs)
	for i := 0; i < 4; i++ {
		writer.WriteU32(p.BuffStates[i])
	}

	// 5. Secondary Stats 부가 정보 직렬화 (Buff 특수 데이터) — 현재는 생략 (필요 시 추후 구현)

	// 6. Secondary Stats 끝 플래그
	writer.WriteU16(0)

	// 7. 직업 (Job)
	writer.WriteU16(p.Character.Class)

	// 8. 외형 Look 직렬화
	p.Character.SerializeLook(writer)

	// 9. 풍선 개수
	writer.WriteU32(p.Balloons)

	// 10. 아이템 이펙트
	writer.WriteU32(p.ItemEffectId)

	// 11. 의자 ID
	writer.WriteU32(p.ChairId)

	// 12. 위치
	writer.Write16(int16(p.Position.X))
	writer.Write16(int16(p.Position.Y))

	// 13. 스탠스
	writer.WriteU8(p.Stance)

	// 14. Foothold (FH)
	writer.WriteU16(0)

	// 15. 펫 정보 (0으로 초기화)
	writer.WriteU8(0)

	// 16. Mount 정보
	writer.WriteU32(p.MountLevel)
	writer.WriteU32(p.MountExp)
	writer.WriteU32(p.MountFatigue)

	// 17. AnnounceBox (0)
	writer.WriteU8(0)

	// 18. 칠판 (Chalkboard)
	if p.Chalkboard != "" {
		writer.WriteU8(1)
		writer.WriteStr16(p.Chalkboard)
	} else {
		writer.WriteU8(0)
	}

	// 19. 반지 정보 직렬화 (Crush, Friendship, Marriage)
	if err := writeRings(writer, p.CrushRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.FriendshipRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.MarriageRings); err != nil {
		return err
	}

	// 20. 반지 끝 표시
	writer.WriteU8(0)

	// 21. Team 정보
	writer.WriteU8(p.Team)

	// 22. 마지막 4바이트 0
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)

	return nil
}

func (s *SpawnPlayer) Deserialize(reader *stream.StreamReader) error {
	return nil
}
