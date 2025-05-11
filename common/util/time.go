package util

import "time"

var KST *time.Location
var TimeFtUtOffset time.Time = FromFileTime(116445060000000000)
var TimeMax time.Time = FromFileTime(150842304000000000)       // -1
var TimeZero time.Time = FromFileTime(94354848000000000)       // -2
var TimePermanent time.Time = FromFileTime(150841440000000000) // -3

func init() {
	var err error
	KST, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic("KST 로케일을 로드할 수 없습니다: " + err.Error())
	}
}

func GetTime(realTimestamp int64) time.Time {
	switch realTimestamp {
	case -1:
		return TimeMax
	case -2:
		return TimeZero
	case -3:
		return TimePermanent
	default:
		return FromFileTime((uint64(realTimestamp) * 10000) + 116445060000000000)
	}
}

func ToFileTime(t time.Time) uint64 {
	const filetimeOffset = uint64(116445060000000000)
	sec := uint64(t.Unix()) * 10_000_000
	nsec := uint64(t.Nanosecond()) / 100
	return filetimeOffset + sec + nsec
}

func FromFileTime(filetime uint64) time.Time {
	const filetimeOffset = uint64(116445060000000000)
	const ticksPerSecond = uint64(10_000_000)

	ft := filetime - filetimeOffset
	sec := int64(ft / ticksPerSecond)
	nsec := int64(ft % ticksPerSecond * 100)
	return time.Unix(sec, nsec)
}
