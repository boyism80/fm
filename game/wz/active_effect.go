package wz

type ActiveEffect struct {
	HP     int
	MP     int
	HPRate int // HP recovery rate (percentage, e.g., 100 = 100%)
	MPRate int // MP recovery rate (percentage, e.g., 100 = 100%)
}
