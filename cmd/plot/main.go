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

// Bind global configuration variables to command-line options
func init() {
	flag.Float64Var(&Z, "Z", Z, "error scale in standard deviations")
	flag.IntVar(&DPI, "dpi", DPI, "DPI for plots")
	flag.StringVar(&PATTERN, "pattern", PATTERN, "file name pattern for generated plots")
	flag.BoolVar(&QUIET, "quiet", QUIET, "make less noise")
}

// describe point for page churn probability
func pagePChurnErr(x, m, s float64) plotter.XYs {
	return plotter.XYs{{x, m}, {x, m + Z*s}, {x, m - Z*s}, {x, m}}
}

// describe mean PPS marker
func ppsMeanErr(mean, std float64) plotter.XYs {
	return plotter.XYs{{mean - Z*std, 0.}, {mean, 0.5}, {mean + Z*std, 0.}}
}

// prepare data for plotting
func parseRecord(record []string) (title string, marker plotter.XYs, pages plotter.XYs) {
	// PPS mean and std
	mean, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Panic(err)
	}
	std, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		log.Panic(err)
	}
	title = fmt.Sprintf("PPS = %.2f±%.2f", mean, Z*std)
	marker = ppsMeanErr(mean, std)

	// per-page churn beliefs
	pages = make(plotter.XYs, (len(record)-3)/2*4) // 4 points per page
	for i := 3; i != len(record); i += 2 {
		// parse beliefs and convert to mean and std
		alpha, err := strconv.ParseFloat(record[i], 64)
		if err != nil {
			log.Panic(err)
		}
		beta, err := strconv.ParseFloat(record[i+1], 64)
		if err != nil {
			log.Panic(err)
		}

		dist := model.Beta{alpha, beta}
		m := dist.Mean()
		s := math.Sqrt(dist.Variance())

		// fill in the page points
		k := (i - 3) / 2
		j := 4 * k
		x := float64(k)
		page := pagePChurnErr(x, m, s)
		copy(pages[j:j+len(page)], page)
	}

	return title, marker, pages
}

// Write to file
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

		// wrap record handling into a function to
		// skip faulty records
		func() {
			iline, _ := strconv.Atoi(record[0])
			defer func() {
				if r := recover(); r != nil {
					log.Printf("skipping "+PATTERN, iline)
				}
			}()

			p, err := plot.New()
			if err != nil {
				log.Panic(err)
			}

			title, marker, pages := parseRecord(record)
			p.Title.Text = fmt.Sprintf("Session %v: %v", iline, title)
			p.X.Label.Text = "Page#"
			plotutil.AddLines(p,
				"P(churn)", pages,
				"PPS", marker)

			writePlot(iline, p)
		}()
	}
}
