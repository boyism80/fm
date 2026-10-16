package game

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewRecordsFromInternalProto(pbs []*internal.Record) *Records {
	r := NewRecords()
	for _, pb := range pbs {
		if pb == nil || pb.GetKey() == "" {
			continue
		}
		record := &Record{
			value:      pb.GetValue(),
			text:       pb.GetText(),
			recordedAt: time.UnixMilli(pb.GetRecordedAtUnixMs()),
		}
		if pb.GetExpiresAtUnixMs() != 0 {
			record.expiresAt = time.UnixMilli(pb.GetExpiresAtUnixMs())
		}
		r.entries[pb.GetKey()] = record
	}
	return r
}

func (r *Records) BuildRPC() []*internal.Record {
	out := make([]*internal.Record, 0, len(r.entries))
	for key := range r.entries {
		record := r.find(key)
		if record == nil {
			continue
		}
		pb := &internal.Record{
			Key:              key,
			Value:            record.value,
			Text:             record.text,
			RecordedAtUnixMs: record.recordedAt.UnixMilli(),
		}
		if !record.expiresAt.IsZero() {
			pb.ExpiresAtUnixMs = record.expiresAt.UnixMilli()
		}
		out = append(out, pb)
	}
	return out
}
