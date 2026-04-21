package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

func PartyDoorFromProto(pb *internal.PartyDoor) *PartyDoor {
	if pb == nil {
		return nil
	}
	return &PartyDoor{
		Town:   pb.GetTown(),
		Target: pb.GetTarget(),
		X:      pb.GetX(),
		Y:      pb.GetY(),
	}
}

func (d *PartyDoor) ToProto() *internal.PartyDoor {
	if d == nil {
		return nil
	}
	return &internal.PartyDoor{
		Town:   d.Town,
		Target: d.Target,
		X:      d.X,
		Y:      d.Y,
	}
}
