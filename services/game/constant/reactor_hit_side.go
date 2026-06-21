package constant

import "fmt"

type ReactorHitSide int32

const (
	ReactorHitAirLeft     ReactorHitSide = 0
	ReactorHitAirRight    ReactorHitSide = 1
	ReactorHitGroundLeft  ReactorHitSide = 2
	ReactorHitGroundRight ReactorHitSide = 3
)

func (s ReactorHitSide) IsLeft() bool {
	return s == ReactorHitAirLeft || s == ReactorHitGroundLeft
}

func (s ReactorHitSide) String() string {
	switch s {
	case ReactorHitAirLeft:
		return "AirLeft"
	case ReactorHitAirRight:
		return "AirRight"
	case ReactorHitGroundLeft:
		return "GroundLeft"
	case ReactorHitGroundRight:
		return "GroundRight"
	default:
		return fmt.Sprintf("ReactorHitSide(%d)", int32(s))
	}
}
