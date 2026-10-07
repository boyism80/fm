package entity

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewMarriageFromInternalProto(pb *internal.Marriage) *Marriage {
	if pb == nil {
		return nil
	}
	m := &Marriage{
		ID:                 pb.GetMarriageId(),
		GroomID:            pb.GetGroomId(),
		BrideID:            pb.GetBrideId(),
		GroomName:          pb.GetGroomName(),
		BrideName:          pb.GetBrideName(),
		GroomItemID:        pb.GetGroomItemId(),
		BrideItemID:        pb.GetBrideItemId(),
		Status:             MarriageStatus(pb.GetStatus()),
		TicketItemID:       pb.GetTicketItemId(),
		GroomWished:        pb.GetGroomWished(),
		BrideWished:        pb.GetBrideWished(),
		GroomWishes:        pb.GetGroomWishes(),
		BrideWishes:        pb.GetBrideWishes(),
		DivorceRequesterID: pb.GetDivorceRequesterId(),
	}
	if ms := pb.GetDivorceRequestedAtUnixMs(); ms > 0 {
		m.DivorceRequestedAt = time.UnixMilli(ms)
	}
	return m
}
