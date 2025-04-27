package dao

type Vec2 struct {
	X, Y int
}

type Object struct {
	ID       int64
	Position Vec2
	Name     string
}
