package model

type Mob struct {
	Life
	DropTable []Item
	Aggro     bool
}
