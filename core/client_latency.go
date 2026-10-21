package core

import "time"

type Latency struct {
	Last     time.Duration
	Smoothed time.Duration
	// Pending is how long the unanswered ping has waited: a stall shows here before any pong arrives.
	Pending time.Duration
}

func (l Latency) Worst() time.Duration {
	return max(l.Last, l.Smoothed, l.Pending)
}

type latencySamples struct {
	pingSentAt time.Time
	last       time.Duration
	smoothed   time.Duration
}

func (c *BaseClient) recordPing(now time.Time) {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()

	c.latency.pingSentAt = now
}

func (c *BaseClient) recordPong(now time.Time) {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()

	if c.latency.pingSentAt.IsZero() {
		return
	}
	rtt := now.Sub(c.latency.pingSentAt)
	c.latency.pingSentAt = time.Time{}
	c.latency.last = rtt
	if c.latency.smoothed == 0 {
		c.latency.smoothed = rtt
		return
	}
	c.latency.smoothed = (7*c.latency.smoothed + rtt) / 8
}

func (c *BaseClient) Latency(now time.Time) Latency {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()

	l := Latency{Last: c.latency.last, Smoothed: c.latency.smoothed}
	if !c.latency.pingSentAt.IsZero() {
		l.Pending = now.Sub(c.latency.pingSentAt)
	}
	return l
}
