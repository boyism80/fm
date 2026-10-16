package game

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

func PartyDoorFromProto(pb *internal.PartyDoor) *PartyDoor {
	if pb == nil {
		return nil
	}
	return &PartyDoor{
		town:   pb.GetTown(),
		target: pb.GetTarget(),
		x:      pb.GetX(),
		y:      pb.GetY(),
	}
}

func (d *PartyDoor) BuildRPC() *internal.PartyDoor {
	if d == nil {
		return nil
	}
	return &internal.PartyDoor{
		Town:   d.town,
		Target: d.target,
		X:      d.x,
		Y:      d.y,
	}
}
