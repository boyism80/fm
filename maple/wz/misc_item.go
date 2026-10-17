package wz

type MiscItem struct {
	*ItemCore
	Exp          int
	QuestID      uint32
	ConsumeItems []uint32
	MaxExp       int
}
