package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type Character struct {
	ID               uint32
	Name             string
	Gender           uint8
	SkinColor        uint8
	Face             uint32
	Hair             uint32
	Pet              uint64
	Level            uint8
	Class            uint16
	Str              uint16
	Dex              uint16
	Int              uint16
	Luk              uint16
	Hp               uint16
	MaxHp            uint16
	Mp               uint16
	MaxMp            uint16
	AbilityPoint     uint16
	Exp              uint32
	Population       uint16
	Map              uint32
	SpawnPoint       uint8
	Rank             uint32
	RankDiff         int32
	ClassRank        uint32
	ClassRankDiff    int32
	Mega             bool
	BaseLooks        map[int8]uint32
	Overlays         map[int8]uint32
	Weapon           uint32
	Position         types.Vector2[int16]
	Stance           uint8
	SkillPoint       uint16
	Inventory        *Inventory
	Skills           []*Skill
	Cooldowns        map[uint32]uint16
	QuestsStarted    []*QuestStatus
	QuestsCompleted  []*QuestStatus
	RegRocks         []uint32
	Rocks            []uint32
	MonsterBookCover uint32
	MonsterBookCards map[uint32]uint32
	RecordExByQuest  map[uint16]string
	Marriage         *Marriage
	Random1          *stream.RandomStream
	Random2          *stream.RandomStream
	Random3          *stream.RandomStream
	BuddyCapacity    uint8
}

func (c *Character) SerializeOverview(writer *stream.StreamWriter) {
	writer.WriteU32(c.ID)
	writer.WriteStaticStr(c.Name, 13)
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteU32(c.Hair)
	writer.WriteU64(c.Pet)
	writer.WriteU8(c.Level)
	writer.WriteU16(c.Class)
	writer.WriteU16(c.Str)
	writer.WriteU16(c.Dex)
	writer.WriteU16(c.Int)
	writer.WriteU16(c.Luk)
	writer.WriteU16(c.Hp)
	writer.WriteU16(c.MaxHp)
	writer.WriteU16(c.Mp)
	writer.WriteU16(c.MaxMp)
	writer.WriteU16(c.AbilityPoint)
	writer.WriteU16(0)
	writer.WriteU32(c.Exp)
	writer.WriteU16(c.Population)
	writer.WriteU32(c.Map)
	writer.WriteU8(c.SpawnPoint)

	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteBoolean(c.Mega)
	writer.WriteU32(c.Hair)

	for parts, itemId := range c.BaseLooks {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	for parts, itemId := range c.Overlays {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	writer.WriteU32(c.Weapon)
	writer.WriteU32(0)

	isRanked := c.Rank > 0
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(c.Rank)
		writer.Write32(c.RankDiff)
		writer.WriteU32(c.ClassRank)
		writer.Write32(c.ClassRankDiff)
	}
}

func (c *Character) DeserializeOverview(reader *stream.StreamReader) {
	c.ID = reader.ReadU32()
	c.Name = reader.ReadStaticStr(13)
	c.Gender = reader.ReadU8()
	c.SkinColor = reader.ReadU8()
	c.Face = reader.ReadU32()
	c.Hair = reader.ReadU32()
	c.Pet = reader.ReadU64()
	c.Level = reader.ReadU8()
	c.Class = reader.ReadU16()
	c.Str = reader.ReadU16()
	c.Dex = reader.ReadU16()
	c.Int = reader.ReadU16()
	c.Luk = reader.ReadU16()
	c.Hp = reader.ReadU16()
	c.MaxHp = reader.ReadU16()
	c.Mp = reader.ReadU16()
	c.MaxMp = reader.ReadU16()
	c.AbilityPoint = reader.ReadU16()
	reader.Skip(2)
	c.Exp = reader.ReadU32()
	c.Population = reader.ReadU16()
	c.Map = reader.ReadU32()
	c.SpawnPoint = reader.ReadU8()

	c.Gender = reader.ReadU8()
	c.SkinColor = reader.ReadU8()
	c.Face = reader.ReadU32()
	c.Mega = reader.ReadBool()
	c.Hair = reader.ReadU32()

	for {
		parts := reader.ReadU8()
		if parts == 0xFF {
			break
		}
		if c.BaseLooks == nil {
			c.BaseLooks = make(map[int8]uint32)
		}
		c.BaseLooks[int8(parts)] = reader.ReadU32()
	}
	for {
		parts := reader.ReadU8()
		if parts == 0xFF {
			break
		}
		if c.Overlays == nil {
			c.Overlays = make(map[int8]uint32)
		}
		c.Overlays[int8(parts)] = reader.ReadU32()
	}

	c.Weapon = reader.ReadU32()
	reader.Skip(4)

	if reader.ReadBool() {
		c.Rank = reader.ReadU32()
		c.RankDiff = reader.Read32()
		c.ClassRank = reader.ReadU32()
		c.ClassRankDiff = reader.Read32()
	}
}

func (c *Character) SerializeLook(writer *stream.StreamWriter) {
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteBoolean(c.Mega)
	writer.WriteU32(c.Hair)

	for parts, itemId := range c.BaseLooks {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	for parts, itemId := range c.Overlays {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	writer.WriteU32(c.Weapon)
	writer.WriteU32(0)
}

func (c *Character) DeserializeLook(reader *stream.StreamReader) {
	c.Gender = reader.ReadU8()
	c.SkinColor = reader.ReadU8()
	c.Face = reader.ReadU32()
	c.Mega = reader.ReadBool()
	c.Hair = reader.ReadU32()

	c.BaseLooks = map[int8]uint32{}
	for parts := int8(reader.ReadU8()); parts != -1; parts = int8(reader.ReadU8()) {
		c.BaseLooks[parts] = reader.ReadU32()
	}

	c.Overlays = map[int8]uint32{}
	for parts := int8(reader.ReadU8()); parts != -1; parts = int8(reader.ReadU8()) {
		c.Overlays[parts] = reader.ReadU32()
	}

	c.Weapon = reader.ReadU32()
	reader.ReadU32()
}
