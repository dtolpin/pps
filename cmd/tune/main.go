// Command tune infers the bandwidth over a sequence of PPS counts.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
    "math"

	"bitbucket.org/dtolpin/pps/model/query"
	"bitbucket.org/dtolpin/pps/infer"
)

func main() {
	bandwidth := flag.Float64("bandwidth", 1000.,
		"initial bandwidth")
    walk := flag.Float64("walk", 100., "standard deviation of random walk")
	total := flag.Int("total", 10,
		"total page count")
	from := flag.Int("from", 0,
		"first row of the tuning set")
    till := flag.Int("till", -1,
        "first row after the tuning set")
    N := flag.Int("N", 1000, "number of MH samples")
	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Create and initialize the model
	m := query.NewModel(*total)
	m.Prior()

	// Read the data
	rdr := csv.NewReader(os.Stdin)
	// assume pps is the last column
	rdr.Read() // skip the header
	// run through sessions
    count := *till - *from
    if count < 0 {
        count = int(*bandwidth)
    }
    counts := make([]int, count)
	for iline := 1; ; iline++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		pps, err := strconv.Atoi(record[len(record)-1])
        counts = append(counts, pps)
		if err != nil {
			log.Printf("line %d: illegal pps %v, skipping",
				iline, record)
			continue
		}
	}

    query := query.NewQuery(*total, *bandwidth, counts)
    proposal := infer.RandomWalk(*walk)
    samples := make(chan float64)
    go infer.MH(query, proposal, *bandwidth, samples)
    sum := 0.
    sum2 := 0.
    for i := 0; i != *N; i++ {
        x := <- samples
        sum += x
        sum2 += x * x
    }
    mean := sum / float64(*N)
    std := math.Sqrt(sum2 / float64(*N) - mean * mean)
    fmt.Printf("bandwidth: %v,%v\n", mean, std)
    log.Printf("bandwidth: %v,%v\n", mean, std)
}
