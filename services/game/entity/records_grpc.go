package entity

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
			Value:      pb.GetValue(),
			Text:       pb.GetText(),
			RecordedAt: time.UnixMilli(pb.GetRecordedAtUnixMs()),
		}
		if pb.GetExpiresAtUnixMs() != 0 {
			record.ExpiresAt = time.UnixMilli(pb.GetExpiresAtUnixMs())
		}
		r.entries[pb.GetKey()] = record
	}
	return r
}

func (r *Records) ToProto() []*internal.Record {
	out := make([]*internal.Record, 0, len(r.entries))
	for key := range r.entries {
		record := r.find(key)
		if record == nil {
			continue
		}
		pb := &internal.Record{
			Key:              key,
			Value:            record.Value,
			Text:             record.Text,
			RecordedAtUnixMs: record.RecordedAt.UnixMilli(),
		}
		if record.ExpiresAt.IsZero() == false {
			pb.ExpiresAtUnixMs = record.ExpiresAt.UnixMilli()
		}
		out = append(out, pb)
	}
	return out
}
