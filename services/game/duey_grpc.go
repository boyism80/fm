package game

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewParcelFromInternalProto(pb *internal.Parcel, gw GameWorld) *Parcel {
	parcel := &Parcel{
		id:      pb.GetParcelId(),
		sender:  pb.GetSenderName(),
		meso:    pb.GetMeso(),
		quick:   pb.GetQuick(),
		message: pb.GetMessage(),
		sentAt:  time.UnixMilli(pb.GetSentAtUnixMs()),
	}
	if pb.GetItem() == nil {
		return parcel
	}
	item, err := NewItemFromInternalProto(pb.GetItem(), gw)
	if err == nil {
		parcel.item = item
	}
	return parcel
}
