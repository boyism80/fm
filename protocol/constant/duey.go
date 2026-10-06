package constant

type DueyMode uint8

const (
	DueyOpenFromArrival DueyMode = 0
	DueyIdentity        DueyMode = 1
	DueySend            DueyMode = 3
	DueyReceive         DueyMode = 5
	DueyDelete          DueyMode = 6
	DueyClose           DueyMode = 8
)

type DueyResult uint8

const (
	DueyResultIdentity          DueyResult = 9
	DueyResultOpen              DueyResult = 10
	DueyResultNotEnoughMeso     DueyResult = 12
	DueyResultInvalid           DueyResult = 13
	DueyResultRecipientNotFound DueyResult = 14
	DueyResultSameAccount       DueyResult = 15
	DueyResultRecipientFull     DueyResult = 16
	DueyResultCannotReceive     DueyResult = 17
	DueyResultOnlyInBox         DueyResult = 18
	DueyResultSent              DueyResult = 19
	DueyResultUnknown           DueyResult = 20
	DueyResultInventoryFull     DueyResult = 22
	DueyResultOnlyHeld          DueyResult = 23
	DueyResultRemoved           DueyResult = 24
	DueyResultAdded             DueyResult = 25
	DueyResultArrival           DueyResult = 26
	DueyResultQuickOpen         DueyResult = 27
	DueyResultArrivals          DueyResult = 28
)

const (
	DueyRemovedDeleted  uint8 = 3
	DueyRemovedReceived uint8 = 4
)
