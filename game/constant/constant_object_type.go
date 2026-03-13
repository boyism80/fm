package constant

// ObjectType is a bitmask for script/Lua entity kind checks.
// Layout: Object=0x0001, Life=Object|0x0010, Character=Object|Life|0x0100, Mob=Object|Life|0x0200,
// Npc=Object|0x0020, Item=Object|0x0040 (0x0030 would include Life bit; Item must not match Life).
type ObjectType uint32

const (
	ObjectTypeObject    ObjectType = 0x0001
	ObjectTypeLife      ObjectType = ObjectTypeObject | 0x0010
	ObjectTypeNpc       ObjectType = ObjectTypeObject | 0x0020
	ObjectTypeItem      ObjectType = ObjectTypeObject | 0x0040
	ObjectTypeCharacter ObjectType = ObjectTypeLife | 0x0100
	ObjectTypeMob       ObjectType = ObjectTypeLife | 0x0200
)

// AllObjectTypeConstants returns name -> value for Lua ObjectType table.
func AllObjectTypeConstants() map[string]ObjectType {
	return map[string]ObjectType{
		"Object":    ObjectTypeObject,
		"Life":      ObjectTypeLife,
		"Character": ObjectTypeCharacter,
		"Mob":       ObjectTypeMob,
		"Npc":       ObjectTypeNpc,
		"Item":      ObjectTypeItem,
	}
}
