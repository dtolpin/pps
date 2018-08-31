package main

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestReadData(t *testing.T) {
	in := `pps
10
8
7
9
5
8
6
10`
	for i, c := range []struct {
		from, till int
		counts     []int
	}{{0, 1000, []int{10, 8, 7, 9, 5, 8, 6, 10}},
		{0, 5, []int{10, 8, 7, 9, 5}},
		{3, 1000, []int{9, 5, 8, 6, 10}},
		{0, 8, []int{10, 8, 7, 9, 5, 8, 6, 10}}} {
		rdr := csv.NewReader(strings.NewReader(in))
		counts := readData(rdr, c.from, c.till)
		if len(counts) != len(c.counts) {
			t.Errorf("readData(%d): wrong number of counts (%v) "+
				"for from=%d, till=%d: got %d, want %d",
				i, counts, c.from, c.till, len(counts), len(c.counts))
		}
		for i := 0; i != len(c.counts); i++ {
			if counts[i] != c.counts[i] {
				t.Errorf("Wrong counts[%d]: got %d, want %d",
					i, counts[i], c.counts[i])
			}
		}
	}
}

func TestInferBandwidth(t *testing.T) {
	for i, c := range []struct {
		total     int
		bandwidth float64
		counts    []int
		walk      float64
		N         int
	}{
		{10, 10000, []int{1, 1, 1, 1, 1, 1, 1, 1}, 10, 1000},
		{10, 10000, []int{1, 10, 1, 10, 1, 10, 1, 10}, 10, 1000},
	} {
		mean, std := inferBandwidth(c.total, c.bandwidth, c.counts,
			c.walk, c.N)
		if mean < std {
			t.Errorf("inferBandwidth(%d) mean (%.3f) should be greated than std (%.3f)",
				i, mean, std)
		}
	}
}
