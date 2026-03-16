package constant

type ObjectType uint32

const (
	ObjectTypeObject    ObjectType = 0x0001
	ObjectTypeLife      ObjectType = ObjectTypeObject | 0x0010
	ObjectTypeNpc       ObjectType = ObjectTypeObject | 0x0020
	ObjectTypeItem      ObjectType = ObjectTypeObject | 0x0040
	ObjectTypeCharacter ObjectType = ObjectTypeLife | 0x0100
	ObjectTypeMob       ObjectType = ObjectTypeLife | 0x0200
	ObjectTypeSummon    ObjectType = ObjectTypeLife | 0x0400
)

func (t ObjectType) Has(other ObjectType) bool {
	return (t & other) == other
}

// AllObjectTypeConstants returns name -> value for Lua ObjectType table.
func AllObjectTypeConstants() map[string]ObjectType {
	return map[string]ObjectType{
		"Object":    ObjectTypeObject,
		"Life":      ObjectTypeLife,
		"Character": ObjectTypeCharacter,
		"Mob":       ObjectTypeMob,
		"Npc":       ObjectTypeNpc,
		"Item":      ObjectTypeItem,
		"Summon":    ObjectTypeSummon,
	}
}
