package constant

type RingActionMode uint8

const (
	RingActionPropose         RingActionMode = 0
	RingActionCancelProposal  RingActionMode = 1
	RingActionAnswer          RingActionMode = 2
	RingActionDropRing        RingActionMode = 3
	RingActionInviteGuest     RingActionMode = 5
	RingActionOpenInvitation  RingActionMode = 6
	RingActionWeddingWishlist RingActionMode = 9
)

type EngageRequestMode uint8

const (
	EngageRequestPropose  EngageRequestMode = 0
	EngageRequestWishlist EngageRequestMode = 9
)

type EngageResult uint8

const (
	EngageResultEngaged             EngageResult = 0x0B
	EngageResultMarried             EngageResult = 0x0C
	EngageResultBroken              EngageResult = 0x0D
	EngageResultDivorced            EngageResult = 0x0E
	EngageResultInvitation          EngageResult = 0x0F
	EngageResultReserved            EngageResult = 0x10
	EngageResultWrongName           EngageResult = 0x12
	EngageResultNotSameMap          EngageResult = 0x13
	EngageResultInventoryFull       EngageResult = 0x14
	EngageResultPartnerInventory    EngageResult = 0x15
	EngageResultSameGender          EngageResult = 0x16
	EngageResultAlreadyEngaged      EngageResult = 0x17
	EngageResultPartnerEngaged      EngageResult = 0x18
	EngageResultAlreadyMarried      EngageResult = 0x19
	EngageResultPartnerMarried      EngageResult = 0x1A
	EngageResultProposalCancelled   EngageResult = 0x1D
	EngageResultDeclined            EngageResult = 0x1E
	EngageResultReservationCanceled EngageResult = 0x1F
	EngageResultCannotCancel        EngageResult = 0x20
	EngageResultInvalidInvitation   EngageResult = 0x22
)

type WeddingPresentMode uint8

const (
	WeddingPresentGive    WeddingPresentMode = 6
	WeddingPresentReceive WeddingPresentMode = 7
	WeddingPresentClose   WeddingPresentMode = 8
)

type WeddingGiftMode uint8

const (
	WeddingGiftOpenGive     WeddingGiftMode = 9
	WeddingGiftOpenReceive  WeddingGiftMode = 10
	WeddingGiftGiven        WeddingGiftMode = 11
	WeddingGiftOneAtATime   WeddingGiftMode = 12
	WeddingGiftGiveFailed   WeddingGiftMode = 13
	WeddingGiftGiveRejected WeddingGiftMode = 14
	WeddingGiftReceived     WeddingGiftMode = 15
	WeddingGiftReceiveFail  WeddingGiftMode = 16
	WeddingGiftOneOfAKind   WeddingGiftMode = 17
)
