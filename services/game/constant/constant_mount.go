package constant

import "time"

const (
	MountMaxLevel        = 30
	MountMaxFatigue      = 100
	MountFatigueInterval = 33 * time.Second
)

var MountExpByLevel = [MountMaxLevel]uint32{0, 6, 25, 50, 105, 134, 196, 254, 263, 315, 367, 430, 543, 587, 679, 725, 897, 1146, 1394, 1701, 2247, 2543, 2898, 3156, 3313, 3584, 3923, 4150, 4305, 4550}

type MountFoodExp struct {
	MinLevel uint32
	Min      int
	Max      int
}

var MountFoodExps = []MountFoodExp{
	{MinLevel: 25, Min: 12, Max: 39},
	{MinLevel: 16, Min: 9, Max: 31},
	{MinLevel: 8, Min: 7, Max: 19},
	{MinLevel: 1, Min: 15, Max: 24},
}
