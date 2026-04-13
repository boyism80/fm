package constant

type ObjectType uint32

const (
	ObjectTypeObject    ObjectType = 0x0001
	ObjectTypeLife      ObjectType = ObjectTypeObject | 0x0010
	ObjectTypeNpc       ObjectType = ObjectTypeObject | 0x0020
	ObjectTypeItem      ObjectType = ObjectTypeObject | 0x0040
	ObjectTypeMist      ObjectType = ObjectTypeObject | 0x0080
	ObjectTypeDoor      ObjectType = ObjectTypeObject | 0x0100
	ObjectTypeCharacter ObjectType = ObjectTypeLife | 0x1000
	ObjectTypeMob       ObjectType = ObjectTypeLife | 0x2000
	ObjectTypeSummon    ObjectType = ObjectTypeLife | 0x4000
)

func (t ObjectType) Has(other ObjectType) bool {
	return (t & other) == other
}

func AllObjectTypeConstants() map[string]ObjectType {
	return map[string]ObjectType{
		"Object":    ObjectTypeObject,
		"Life":      ObjectTypeLife,
		"Character": ObjectTypeCharacter,
		"Mob":       ObjectTypeMob,
		"Npc":       ObjectTypeNpc,
		"Item":      ObjectTypeItem,
		"Summon":    ObjectTypeSummon,
		"Mist":      ObjectTypeMist,
		"Door":      ObjectTypeDoor,
	}
}
