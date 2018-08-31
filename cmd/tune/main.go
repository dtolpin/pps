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

	"image"
	"image/gif"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"bitbucket.org/dtolpin/pps/infer"
	"bitbucket.org/dtolpin/pps/model/query"
)

// Be random.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// The command line
var BANDWIDTH float64 = 1000
var WALK float64 = 10
var TOTAL int = 20
var FROM int = 0
var TILL int = -1
var N int = 1000
var Z float64 = 2
var PLOT string = ""
var DPI int = 120

func init() {
	flag.Float64Var(&BANDWIDTH, "bandwidth", BANDWIDTH,
		"initial bandwidth")
	flag.Float64Var(&WALK, "walk", WALK, "standard deviation of random walk")
	flag.IntVar(&TOTAL, "total", TOTAL,
		"total page count")
	flag.IntVar(&FROM, "from", FROM,
		"first row of the tuning set")
	flag.IntVar(&TILL, "till", TILL,
		"first row after the tuning set")
	flag.IntVar(&N, "N", N, "number of MH samples")
	flag.Float64Var(&Z, "Z", Z, "Z score for confidence range in output")
	flag.StringVar(&PLOT, "plot", PLOT, "file to store the plot")
	flag.IntVar(&DPI, "dpi", DPI, "DPI of the plot")
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

// Function drawBandwidth draws the inferred bandwidth as a histogram.
func drawBandwidth(dist []float64, mean, std float64) {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.Title.Text = fmt.Sprintf("Bandwidth: %.0f ± %.0f", mean, Z*std)
	p.X.Label.Text = "bandwidth"
	h, err := plotter.NewHist(plotter.Values(dist), 16)
	if err != nil {
		panic(err)
	}
	// Normalize the area under the histogram to
	// sum to one.
	h.Normalize(1)
	p.Add(h)

	img := image.NewRGBA(image.Rect(0, 0, 7*DPI, 3*DPI))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	f, err := os.Create(PLOT)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}

	if err := gif.Encode(f, img, nil); err != nil {
		log.Panic(err)
	}
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
	dist := make([]float64, N)

	// Burn
	for i := 0; i != N; i++ {
		<-samples
	}

	// Collect after burn-in
	for i := 0; i != N; i++ {
		dist[i] = <-samples
		sum += dist[i]
		sum2 += dist[i] * dist[i]
	}
	mean = sum / float64(N)
	std = math.Sqrt(sum2/float64(N) - mean*mean)

	if PLOT != "" {
		drawBandwidth(dist, mean, std)
	}

	return mean, std
}

func main() {
	flag.Parse()
	if TILL == -1 {
		TILL = math.MaxInt32
	}

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Read the PPS data
	counts := readData(csv.NewReader(os.Stdin), FROM, TILL)

	// Infer the bandwidth
	mean, std := inferBandwidth(TOTAL, BANDWIDTH, counts,
		WALK, N)
	fmt.Printf("Bandwidth: %.0f ± %.0f\n", mean, Z*std)
}
