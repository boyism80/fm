package types

type ObjectType uint32

const (
	OBJECT_TYPE_NONE           ObjectType = 0
	OBJECT_TYPE_NPC            ObjectType = 1 << 0
	OBJECT_TYPE_MONSTER        ObjectType = 1 << 1
	OBJECT_TYPE_ITEM           ObjectType = 1 << 2
	OBJECT_TYPE_PLAYER         ObjectType = 1 << 3
	OBJECT_TYPE_DOOR           ObjectType = 1 << 4
	OBJECT_TYPE_SUMMON         ObjectType = 1 << 5
	OBJECT_TYPE_SHOP           ObjectType = 1 << 6
	OBJECT_TYPE_MIST           ObjectType = 1 << 7
	OBJECT_TYPE_REACTOR        ObjectType = 1 << 8
	OBJECT_TYPE_MESSAGE_BOX    ObjectType = 1 << 9
	OBJECT_TYPE_HIRED_MERCHANT ObjectType = 1 << 10

	// 조합??
	OBJECT_TYPE_OBJECT ObjectType = OBJECT_TYPE_MONSTER | OBJECT_TYPE_PLAYER | OBJECT_TYPE_NPC | OBJECT_TYPE_SUMMON | OBJECT_TYPE_ITEM | OBJECT_TYPE_DOOR | OBJECT_TYPE_REACTOR | OBJECT_TYPE_SHOP | OBJECT_TYPE_MESSAGE_BOX | OBJECT_TYPE_HIRED_MERCHANT
	OBJECT_TYPE_LIFE   ObjectType = OBJECT_TYPE_MONSTER | OBJECT_TYPE_PLAYER
)
