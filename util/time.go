package util

import "time"

var KST *time.Location

const fileTimeOffset = uint64(116445060000000000)
const fileTimeTicksPerSecond = uint64(10_000_000)

const FileTimeZero = uint64(94354848000000000)

var TimeMax time.Time = FromFileTime(150842304000000000)

func init() {
	var err error
	KST, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic("KST 로케일을 로드할 수 없습니다: " + err.Error())
	}
}

func ToFileTime(t time.Time) uint64 {
	return uint64(t.UnixNano()/100) + fileTimeOffset
}

func FromFileTime(ft uint64) time.Time {
	return time.Unix(int64((ft-fileTimeOffset)/fileTimeTicksPerSecond), int64(((ft-fileTimeOffset)%fileTimeTicksPerSecond)*100))
}
