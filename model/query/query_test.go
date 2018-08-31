package query

import (
	"math"
	"testing"
)

const epsilon = 1E-6

func TestQueryObserve(t *testing.T) {
	for i, c := range []struct {
		total             int
		sample, bandwidth float64
		counts            []int
		logp              float64
	}{
		{1, 1000., 1000., []int{1, 1}, -1. + math.Log(0.5) + math.Log(0.75)},
		{2, 1000., 2000., []int{1, 1}, -0.5 + math.Log(0.5) + math.Log(0.75)},
		{2, 1000., 500., []int{1, 2}, -2. + math.Log(0.5) + math.Log(0.25)},
		{4, 1000., 1000., []int{2, 2}, -1. + math.Log(0.5) + math.Log(0.5) +
			math.Log(0.75) + math.Log(0.75)},
		{1, -10., 1000., []int{}, -math.MaxFloat64},
	} {
		q := NewQuery(c.total, c.bandwidth, c.counts)
		logp := q.Observe(c.sample)
		if math.Abs(logp-c.logp) > epsilon {
			t.Errorf("Query.Observe(%d): wrong logp: got %.3g want %.3g",
				i, logp, c.logp)
		}
	}
}
