package main

import (
    "flag"
	"fmt"
	"io"
	"log"
	"os"

	"encoding/csv"
	"math"
	"strconv"

	"image"
	"image/gif"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

    "bitbucket.org/dtolpin/pps/model"
)

var Z float64 = 2.
var DPI int = 120
var PATTERN string = "pps-%06v.gif"
var QUIET bool = false
func init () {
    flag.Float64Var(&Z, "Z", Z, "error scale in standard deviations")
    flag.IntVar(&DPI, "dpi", DPI, "DPI for plots")
    flag.StringVar(&PATTERN, "pattern", PATTERN, "file name pattern for generated plots")
    flag.BoolVar(&QUIET, "quiet", QUIET, "make less noise")
}

// define point with error bar
func pterr(pts plotter.XYs, i int, m, s float64) {
	k := (i - 3) / 2
	j := 4 * k
	x := float64(k)
	pts[j].X = x
	pts[j].Y = m
	pts[j+1].X = x
	pts[j+1].Y = m + Z*s
	pts[j+2].X = x
	pts[j+2].Y = m - Z*s
	pts[j+3].X = x
	pts[j+3].Y = m
}

// define PPS mean/std marker
func meanstd(mean, std float64) plotter.XYs {
    return plotter.XYs{
        {mean - Z * std, 0.},
        {mean, 0.5},
        {mean + Z * std, 0.}}
}

// define points and mean for record
func points(record []string) (float64, float64, plotter.XYs) {
	mean, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Panic(err)
	}
	std, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		fmt.Println(err)
	}
	pts := make(plotter.XYs, (len(record) - 3)/2*4)
	for i := 3; i != len(record); i += 2 {
		alpha, err := strconv.ParseFloat(record[i], 64)
		if err != nil {
			fmt.Println(err)
		}
		beta, err := strconv.ParseFloat(record[i+1], 64)
		if err != nil {
			fmt.Println(err)
		}
        dist := model.Beta{alpha, beta}
		m := dist.Mean()
		s := math.Sqrt(dist.Variance())
		pterr(pts, i, m, s)
	}
	return mean, std, pts
}

// Write to file.
func writePlot(iline int, p *plot.Plot) {
    img := image.NewRGBA(image.Rect(0, 0, 7*DPI, 3*DPI))
    c := vgimg.NewWith(vgimg.UseImage(img))
    p.Draw(draw.New(c))

    fname := fmt.Sprintf(PATTERN, iline)
    if !QUIET {
        log.Println(fname)
    }

    f, err := os.Create(fmt.Sprintf(PATTERN, iline))
    defer f.Close()
    if err != nil {
        log.Panic(err)
    }

    if err := gif.Encode(f, img, nil); err != nil {
        log.Panic(err)
    }
}

func main() {
    flag.Parse()

	rdr := csv.NewReader(os.Stdin)
    rdr.Read() // skip the header
	for i := 0; ; i++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		p, err := plot.New()
		mean, std, pts := points(record)
        iline, _ := strconv.Atoi(record[0])

		p.Title.Text = fmt.Sprintf("session %v, PPS = %.2f±%.2f", iline, mean, Z*std)
		p.X.Label.Text = "page#"

		plotutil.AddLines(p,
			"P(Churn)", pts,
			"mean(PPS)", meanstd(mean, std))

        writePlot(iline, p)
	}
}
