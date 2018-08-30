package main

import (
    "testing"
	"encoding/csv"
	"strings"
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
    for _, c := range []struct{
        from, till int
        counts []int
    }{{0, 1000, []int{10, 8, 7, 9, 5, 8, 6, 10}},
      {0, 5, []int{10, 8, 7, 9, 5}},
      {3, 1000, []int{9, 5, 8, 6, 10}},
      {0, 8, []int{10, 8, 7, 9, 5, 8, 6, 10}}} {
        rdr := csv.NewReader(strings.NewReader(in))
        counts := readData(rdr, c.from, c.till)
        if(len(counts) != len(c.counts)) {
            t.Errorf("Wrong number of counts (%v) for from=%d, till=%d: got %d, want %d",
                counts, c.from, c.till, len(counts), len(c.counts))
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
}
