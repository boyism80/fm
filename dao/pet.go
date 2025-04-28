package dao

type Pet struct {
	Name        string
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	PetItemId   uint32
	SecondsLeft uint32
	Expiration  int64
}
