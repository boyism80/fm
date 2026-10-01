package entity

import "testing"

func TestGuildLevel(t *testing.T) {
	cases := []struct {
		gp    uint32
		level uint32
	}{
		{0, 1},
		{19_999, 1},
		{20_000, 2},
		{159_999, 2},
		{160_000, 3},
		{14_579_999, 9},
		{14_580_000, 10},
		{4_000_000_000, 10},
	}
	for _, c := range cases {
		g := &Guild{GP: c.gp}
		if got := g.Level(); got != c.level {
			t.Errorf("GP %d: level %d, want %d", c.gp, got, c.level)
		}
	}
}
