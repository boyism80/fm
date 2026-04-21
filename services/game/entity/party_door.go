package entity

// PartyDoor is in-game party door state (not the gRPC message type).
type PartyDoor struct {
	Town   uint32
	Target uint32
	X      int32
	Y      int32
}

func (d *PartyDoor) Clone() *PartyDoor {
	if d == nil {
		return nil
	}
	return &PartyDoor{
		Town:   d.Town,
		Target: d.Target,
		X:      d.X,
		Y:      d.Y,
	}
}
