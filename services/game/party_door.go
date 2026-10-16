package game

// PartyDoor is in-game party door state (not the gRPC message type).
type PartyDoor struct {
	town   uint32
	target uint32
	x      int32
	y      int32
}

func (d *PartyDoor) Clone() *PartyDoor {
	if d == nil {
		return nil
	}
	return &PartyDoor{
		town:   d.town,
		target: d.target,
		x:      d.x,
		y:      d.y,
	}
}
