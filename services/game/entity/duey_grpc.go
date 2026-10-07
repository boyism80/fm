package entity

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewParcelFromInternalProto(pb *internal.Parcel, gw GameWorld) *Parcel {
	parcel := &Parcel{
		ID:      pb.GetParcelId(),
		Sender:  pb.GetSenderName(),
		Meso:    pb.GetMeso(),
		Quick:   pb.GetQuick(),
		Message: pb.GetMessage(),
		SentAt:  time.UnixMilli(pb.GetSentAtUnixMs()),
	}
	if pb.GetItem() == nil {
		return parcel
	}
	item, err := NewItemFromInternalProto(pb.GetItem(), gw)
	if err == nil {
		parcel.Item = item
	}
	return parcel
}
