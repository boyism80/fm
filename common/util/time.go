// Package util provides common utility functions for the MapleStory private server.
// This file contains time-related utilities for handling MapleStory's file time format
// and Korean Standard Time operations.
package util

import "time"

// KST represents the Korean Standard Time location.
var KST *time.Location

// MapleStory time constants for special timestamp values
var TimeFtUtOffset time.Time = FromFileTime(116445060000000000)
var TimeMax time.Time = FromFileTime(150842304000000000)       // -1: never expires
var TimeZero time.Time = FromFileTime(94354848000000000)       // -2: no expiration
var TimePermanent time.Time = FromFileTime(150841440000000000) // -3: permanent

// init initializes the Korean Standard Time location.
// Panics if KST location cannot be loaded, as this is critical for server operation.
func init() {
	var err error
	KST, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic("KST 로케일을 로드할 수 없습니다: " + err.Error())
	}
}

// GetTime converts MapleStory timestamp to Go time.Time.
// Special values: -1 (max), -2 (zero), -3 (permanent)
//
// Parameters:
//   - realTimestamp: The timestamp from MapleStory data
//
// Returns:
//   - time.Time: Converted Go time value
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

// ToFileTime converts Go time.Time to Windows FILETIME format.
// FILETIME represents the number of 100-nanosecond intervals since January 1, 1601 UTC.
// This is used when sending time data to the MapleStory client.
//
// Parameters:
//   - t: The Go time.Time to convert
//
// Returns:
//   - uint64: Time in FILETIME format (100-nanosecond intervals since 1601-01-01 UTC)
func ToFileTime(t time.Time) uint64 {
	const FILETIME_OFFSET = uint64(116445060000000000)
	return uint64(t.UnixNano()/100) + FILETIME_OFFSET
}

// FromFileTime converts Windows FILETIME to Go time.Time.
// FILETIME represents the number of 100-nanosecond intervals since January 1, 1601 UTC.
// This is used when reading time data from MapleStory files.
//
// Parameters:
//   - ft: FILETIME value (100-nanosecond intervals since 1601-01-01 UTC)
//
// Returns:
//   - time.Time: Converted Go time value
func FromFileTime(ft uint64) time.Time {
	const FILETIME_OFFSET = uint64(116445060000000000)
	const TICKS_PER_SECOND = uint64(10_000_000)
	return time.Unix(int64((ft-FILETIME_OFFSET)/TICKS_PER_SECOND), int64(((ft-FILETIME_OFFSET)%TICKS_PER_SECOND)*100))
}
