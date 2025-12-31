package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/protocol/dto"
)

type SpawnPlayer struct {
	Character         *dto.Character
	GuildName         string    // Guild name (empty if none)
	GuildLogoBG       uint16    // Guild logo background shape
	GuildLogoBGColor  uint8     // 길드 ?블??배경 ??
	GuildLogo         uint16    // Guild logo background shape
	GuildLogoColor    uint8     // 길드 ?블????
	BuffStates        [4]uint32 // Buff bitmask
	Diseases          [4]uint32 // Disease bitmask
	SpeedBuff         uint8     // Buff speed (optional)
	ComboCount        uint8     // Combo count (optional)
	WKChargeSkillId   uint32    // WK charge skill ID (optional)
	MorphId           uint16    // Morph ID (optional)
	SpiritClawSkillId uint32    // Spirit claw skill ID (optional)
	ItemEffectId      uint32    // 캐릭???용 ?이???펙??
	ChairId           uint32    // Chair ID being sat on
	Balloons          uint32    // Balloon count
	MountLevel        uint32    // Mount level
	MountExp          uint32    // ?것 경험?
	MountFatigue      uint32    // ?것 ?로??
	Chalkboard        string    // Chalkboard text (empty if none)
	Team              uint8     // Team (0/1)
	CrushRings        []*dto.Ring
	FriendshipRings   []*dto.Ring
	MarriageRings     []*dto.Ring
}

func writeRings(writer *stream.StreamWriter, rings []*dto.Ring) error {
	writer.WriteU8(uint8(len(rings))) // Ring count
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		writer.WriteU64(ring.RingId)       // Unique ring ID
		writer.WriteU64(ring.PartnerId)    // Partner character ID
		writer.WriteU64(ring.RingUniqueId) // Ring unique UID (0 if not equipped)
	}
	return nil
}

func (p *SpawnPlayer) Serialize(writer *stream.StreamWriter) error {
	if p.Character == nil {
		return nil
	}

	// Partner character ID
	writer.WriteU32(p.Character.ID)

	// 2. Character name (MapleAsciiString)
	writer.WriteStr16(p.Character.Name)

	// 3. Guild info
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
	for i := range 4 {
		writer.WriteU32(p.BuffStates[i])
	}

	// 5. Secondary Stats additional info serialization (Buff count item count etc. - simplified, implement later)

	// 6. Secondary Stats ???래?
	writer.WriteU16(0)

	// 7. Job
	writer.WriteU16(p.Character.Class)

	// 8. ?형 Look 직렬??
	p.Character.SerializeLook(writer)

	// Balloon count
	writer.WriteU32(p.Balloons)

	// 10. ?이???펙??
	writer.WriteU32(p.ItemEffectId)

	// 11. Chair ID
	writer.WriteU32(p.ChairId)

	// 12. Position
	writer.Write16(p.Character.Position.X)
	writer.Write16(p.Character.Position.Y)

	// 13. ?탠??
	writer.WriteU8(p.Character.Stance)

	// 14. Foothold (FH)
	writer.WriteU16(0)

	// 15. ???보 (0?로 초기??
	writer.WriteU8(0)

	// 16. Mount info
	writer.WriteU32(p.MountLevel)
	writer.WriteU32(p.MountExp)
	writer.WriteU32(p.MountFatigue)

	// 17. AnnounceBox (0)
	writer.WriteU8(0)

	// 18. Chalkboard
	if p.Chalkboard != "" {
		writer.WriteU8(1)
		writer.WriteStr16(p.Chalkboard)
	} else {
		writer.WriteU8(0)
	}

	// 19. Ring info serialization (Crush, Friendship, Marriage)
	if err := writeRings(writer, p.CrushRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.FriendshipRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.MarriageRings); err != nil {
		return err
	}

	// 20. Ring unknown
	writer.WriteU8(0)

	// 21. Team info
	writer.WriteU8(p.Team)

	// 22. 마??4바이??0
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)

	return nil
}

func (s *SpawnPlayer) Opcode() uint16 {
	return 0x6E // SpawnPlayer opcode
}

func (s *SpawnPlayer) Deserialize(reader *stream.StreamReader) error {
	return nil
}
