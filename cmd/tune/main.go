// Command tune infers the bandwidth over a sequence of PPS counts.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"bitbucket.org/dtolpin/pps/infer"
	"bitbucket.org/dtolpin/pps/model/query"
)

// Be random.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func readData(rdr *csv.Reader, from, till int) []int {
	rdr.Read() // skip the header

	counts := make([]int, 0)
	for iline := 0; iline != till; iline++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if iline >= from {
			pps, err := strconv.Atoi(record[len(record)-1])
			if err != nil {
				log.Printf("line %d: illegal pps %v, skipping",
					iline, record)
				continue
			}
			counts = append(counts, pps)
		}
	}
	return counts
}

// Function inferBandwidth invokes approximate inference to infer
// the bandwidth. Returns mean and variance of the bandwidth.
func inferBandwidth(total int, bandwidth float64, counts []int,
	walk float64, N int) (mean, std float64) {
	query := query.NewQuery(total, bandwidth, counts)
	proposal := infer.RandomWalk(walk)
	samples := make(chan float64)
	go infer.MH(query, proposal, bandwidth, samples)
	sum := 0.
	sum2 := 0.
	for i := 0; i != N; i++ {
		x := <-samples
		sum += x
		sum2 += x * x
	}
	mean = sum / float64(N)
	std = math.Sqrt(sum2/float64(N) - mean*mean)
	return mean, std
}

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
	Z := flag.Float64("Z", 2., "Z score for confidence range in output")
	flag.Parse()
	if *till == -1 {
		*till = math.MaxInt32
	}

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Read the PPS data
	counts := readData(csv.NewReader(os.Stdin), *from, *till)

	// Infer the bandwidth
	mean, std := inferBandwidth(*total, *bandwidth, counts,
		*walk, *N)
	fmt.Printf("Bandwidth: %.0f ± %.0f\n", mean, *Z*std)
}
