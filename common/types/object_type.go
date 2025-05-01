package types

type ObjectType uint32

const (
	ObjectTypeNone          ObjectType = 0
	ObjectTypeNPC           ObjectType = 1 << 0
	ObjectTypeMonster       ObjectType = 1 << 1
	ObjectTypeItem          ObjectType = 1 << 2
	ObjectTypePlayer        ObjectType = 1 << 3
	ObjectTypeDoor          ObjectType = 1 << 4
	ObjectTypeSummon        ObjectType = 1 << 5
	ObjectTypeShop          ObjectType = 1 << 6
	ObjectTypeMist          ObjectType = 1 << 7
	ObjectTypeReactor       ObjectType = 1 << 8
	ObjectTypeMessageBox    ObjectType = 1 << 9
	ObjectTypeHiredMerchant ObjectType = 1 << 10

	// 조합형
	ObjectTypeObject ObjectType = ObjectTypeMonster | ObjectTypePlayer | ObjectTypeNPC | ObjectTypeSummon | ObjectTypeItem | ObjectTypeDoor | ObjectTypeReactor | ObjectTypeShop | ObjectTypeMessageBox | ObjectTypeHiredMerchant
	ObjectTypeLife   ObjectType = ObjectTypeMonster | ObjectTypePlayer
)
