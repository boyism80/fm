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
		r.entries[pb.GetKey()] = &Record{
			Value:     pb.GetValue(),
			Text:      pb.GetText(),
			Period:    RecordPeriod(pb.GetPeriod()),
			UpdatedAt: time.UnixMilli(pb.GetUpdatedAtUnixMs()),
		}
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
		out = append(out, &internal.Record{
			Key:             key,
			Value:           record.Value,
			Text:            record.Text,
			Period:          uint32(record.Period),
			UpdatedAtUnixMs: record.UpdatedAt.UnixMilli(),
		})
	}
	return out
}
