package constant

type BuddyMode uint8

const (
	BuddyAdd    BuddyMode = 1
	BuddyAccept BuddyMode = 2
	BuddyDelete BuddyMode = 3
)

const BuddyDefaultGroup = "그룹 미지정"

const DefaultBuddyCapacity uint8 = 20

type BuddyS2CSubOpcode uint8

const (
	BuddyS2CAddRequest     BuddyS2CSubOpcode = 9
	BuddyS2CChannelUpdate  BuddyS2CSubOpcode = 0x14
	BuddyS2CCapacityUpdate BuddyS2CSubOpcode = 0x15
)

type BuddyListSyncAction uint8

const (
	BuddyListSyncLogin  BuddyListSyncAction = 7
	BuddyListSyncUpdate BuddyListSyncAction = 10
	BuddyListSyncDelete BuddyListSyncAction = 18
)

type BuddyStatusCode uint8

const (
	BuddyStatusGeneric    BuddyStatusCode = 11
	BuddyStatusTargetFull BuddyStatusCode = 12
	BuddyStatusNotFound   BuddyStatusCode = 15
)

func BuddyStatusForInternalError(errorCode int32) BuddyStatusCode {
	switch errorCode {
	case 2:
		return BuddyStatusNotFound
	case 4:
		return BuddyStatusTargetFull
	default:
		return BuddyStatusGeneric
	}
}
