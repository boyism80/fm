package server

import (
	"sync"
	"time"
)

const (
	uniqueIDEpoch    = 1767225600
	uniqueIDNodeBits = 10
	uniqueIDSeqBits  = 12
	uniqueIDSeqLimit = 1 << uniqueIDSeqBits
)

type uniqueIDs struct {
	mu     sync.Mutex
	node   uint64
	second int64
	seq    uint64
}

func (u *uniqueIDs) next() uint64 {
	u.mu.Lock()
	defer u.mu.Unlock()

	now := max(time.Now().Unix()-uniqueIDEpoch, u.second)
	if now == u.second {
		u.seq++
		if u.seq == uniqueIDSeqLimit {
			now++
			u.seq = 0
		}
	} else {
		u.seq = 0
	}
	u.second = now
	return uint64(now)<<(uniqueIDNodeBits+uniqueIDSeqBits) | u.node<<uniqueIDSeqBits | u.seq
}
