package entity

import (
	"time"

	"github.com/boyism80/fm/core/clock"
)

func (qp *Quest) IsDeadlineExpired() bool {
	if qp == nil || qp.Deadline.IsZero() {
		return false
	}
	return !clock.Now().Before(qp.Deadline)
}

func (qp *Quest) SetDeadline(t time.Time) {
	if qp == nil {
		return
	}
	qp.Deadline = t
}

func (qp *Quest) SetDeadlineAfter(d time.Duration) {
	if qp == nil || d <= 0 {
		return
	}
	qp.Deadline = clock.Now().Add(d)
}

func (qp *Quest) ClearDeadline() {
	if qp == nil {
		return
	}
	qp.Deadline = time.Time{}
}
