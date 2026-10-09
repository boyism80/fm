package uniqueid

import (
	"fmt"
	"sync"
	"time"
)

const (
	epochMs      = 1767225600000
	nodeBits     = 10
	seqBits      = 12
	seqLimit     = 1 << seqBits
	worldNodeMax = 1 << (nodeBits - 5)
	serverMax    = 1 << 4
)

type Generator struct {
	mu   sync.Mutex
	node uint64
	ms   int64
	seq  uint64
}

func NewForChannel(worldID uint32, channelID uint32) *Generator {
	if worldID >= worldNodeMax || channelID >= serverMax {
		panic(fmt.Sprintf("uniqueid: world %d channel %d out of range", worldID, channelID))
	}
	return &Generator{node: uint64(worldID)<<5 | uint64(channelID)}
}

func NewForCashShop(worldID uint32, cashShopID uint32) *Generator {
	if worldID >= worldNodeMax || cashShopID >= serverMax {
		panic(fmt.Sprintf("uniqueid: world %d cash shop %d out of range", worldID, cashShopID))
	}
	return &Generator{node: uint64(worldID)<<5 | serverMax | uint64(cashShopID)}
}

func (u *Generator) Next() uint64 {
	u.mu.Lock()
	defer u.mu.Unlock()

	now := max(time.Now().UnixMilli()-epochMs, u.ms)
	if now == u.ms {
		u.seq++
		if u.seq == seqLimit {
			now++
			u.seq = 0
		}
	} else {
		u.seq = 0
	}
	u.ms = now
	return uint64(now)<<(nodeBits+seqBits) | u.node<<seqBits | u.seq
}
