package query

import (
	"math"
	"testing"
)

const epsilon = 1E-6

func TestQueryObserve(t *testing.T) {
	bandwidth := 1000.
	for _, c := range []struct {
		total     int
		bandwidth float64
		counts    []int
		logp      float64
	}{{1, 1000., []int{1, 1}, -1. + math.Log(0.5) + math.Log(0.75)},
		{2, 2000., []int{1, 1}, -0.5 + math.Log(0.5) + math.Log(0.75)},
		{2, 500., []int{1, 2}, -2. + math.Log(0.5) + math.Log(0.25)},
		{4, 1000., []int{2, 2}, -1. + math.Log(0.5) + math.Log(0.5) +
			math.Log(0.75) + math.Log(0.75)}} {
		q := NewQuery(c.total, c.bandwidth, c.counts)
		logp := q.Observe(bandwidth)
		if math.Abs(logp-c.logp) > epsilon {
			t.Errorf("wrong logp for total=%d, bandwidth=%.3g, counts=%v: "+
				"got %.3g want %.3g",
				c.total, c.bandwidth, c.counts, logp, c.logp)
		}
	}
}
