package wz

// Healer is the healing object a field lowers while players pull its pulley; it rises back on its own.
type Healer struct {
	YMin    int16
	YMax    int16
	Fall    int16
	Rise    int16
	HealMin int
	HealMax int
}
