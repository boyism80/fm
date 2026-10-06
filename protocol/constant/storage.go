package constant

type StorageMode uint8

const (
	StorageTakeOut StorageMode = 4
	StorageStore   StorageMode = 5
	StorageArrange StorageMode = 6
	StorageMeso    StorageMode = 7
	StorageClose   StorageMode = 8
)

type StorageResult uint8

const (
	StorageResultTakeOut       StorageResult = 0x09
	StorageResultInventoryFull StorageResult = 0x0A
	StorageResultNotEnoughMeso StorageResult = 0x0B
	StorageResultStore         StorageResult = 0x0D
	StorageResultArrange       StorageResult = 0x0F
	StorageResultFull          StorageResult = 0x11
	StorageResultMeso          StorageResult = 0x13
	StorageResultOpen          StorageResult = 0x16
)
