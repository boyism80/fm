package game

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewMarriageFromInternalProto(pb *internal.Marriage) *Marriage {
	if pb == nil {
		return nil
	}
	m := &Marriage{
		id:                 pb.GetMarriageId(),
		groomID:            pb.GetGroomId(),
		brideID:            pb.GetBrideId(),
		groomName:          pb.GetGroomName(),
		brideName:          pb.GetBrideName(),
		groomItemID:        pb.GetGroomItemId(),
		brideItemID:        pb.GetBrideItemId(),
		status:             MarriageStatus(pb.GetStatus()),
		ticketItemID:       pb.GetTicketItemId(),
		groomWished:        pb.GetGroomWished(),
		brideWished:        pb.GetBrideWished(),
		groomWishes:        pb.GetGroomWishes(),
		brideWishes:        pb.GetBrideWishes(),
		divorceRequesterID: pb.GetDivorceRequesterId(),
	}
	if ms := pb.GetDivorceRequestedAtUnixMs(); ms > 0 {
		m.divorceRequestedAt = time.UnixMilli(ms)
	}
	return m
}
