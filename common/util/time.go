package util

const (
	TimeFtUtOffset uint64 = 116445060000000000 // KST
	TimeMax        uint64 = 150842304000000000 // 00 80 05 BB 46 E6 17 02
	TimeZero       uint64 = 94354848000000000  // 00 40 E0 FD 3B 37 4F 01
	TimePermanent  uint64 = 150841440000000000 // 00 C0 9B 90 7D E5 17 02
)

// GetTime converts a real timestamp to a time format based on specific conditions
func GetTime(realTimestamp int64) uint64 {
	switch realTimestamp {
	case -1:
		return TimeMax
	case -2:
		return TimeZero
	case -3:
		return TimePermanent
	default:
		return (uint64(realTimestamp) * 10000) + TimeFtUtOffset
	}
}

func GetKoreanTimestamp(realTimestamp int64) uint64 {
	return GetTime(realTimestamp)
}
