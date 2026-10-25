package util

import "time"

// fileTimeEpoch is 1970-01-01 00:00 UTC as a FILETIME.
const fileTimeEpoch = int64(116444736000000000)
const fileTimeTicksPerSecond = int64(10_000_000)

const FileTimeZero = uint64(94354848000000000)

var TimeMax time.Time = FromFileTime(150842304000000000)

// The client reads a FILETIME as the server's local wall time, so the zone offset at t is added.
func ToFileTime(t time.Time) uint64 {
	_, offset := t.In(time.Local).Zone()
	return uint64(t.UnixNano()/100 + int64(offset)*fileTimeTicksPerSecond + fileTimeEpoch)
}

func FromFileTime(ft uint64) time.Time {
	ticks := int64(ft) - fileTimeEpoch
	wall := time.Unix(ticks/fileTimeTicksPerSecond, (ticks%fileTimeTicksPerSecond)*100)
	_, offset := wall.In(time.Local).Zone()
	return wall.Add(-time.Duration(offset) * time.Second)
}
