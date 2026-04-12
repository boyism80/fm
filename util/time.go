package util

import "time"

var KST *time.Location

var TimeFtUtOffset time.Time = FromFileTime(116445060000000000)
var TimeMax time.Time = FromFileTime(150842304000000000)
var TimeZero time.Time = FromFileTime(94354848000000000)
var TimePermanent time.Time = FromFileTime(150841440000000000)

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
	const FILETIME_OFFSET = uint64(116445060000000000)
	return uint64(t.UnixNano()/100) + FILETIME_OFFSET
}

func FromFileTime(ft uint64) time.Time {
	const FILETIME_OFFSET = uint64(116445060000000000)
	const TICKS_PER_SECOND = uint64(10_000_000)
	return time.Unix(int64((ft-FILETIME_OFFSET)/TICKS_PER_SECOND), int64(((ft-FILETIME_OFFSET)%TICKS_PER_SECOND)*100))
}
