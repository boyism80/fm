package model

type Character struct {
	Life
	PlayerID int64
	Class    string
	Exp      int
}
