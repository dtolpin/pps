// Command scan runs over pps sequence (stored as a CSV file, last column is the pps)
// and computes updated beliefs after each session. The output is enumerated CSV, where
// each line is a flattened Beliefs vector.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"bitbucket.org/dtolpin/pps/model"
)

// Command line
var BANDWIDTH float64 = 1000
var TOTAL int = 30
var THIN int = 20
var FLOATFMT = "%.3g"

func init() {
	flag.Float64Var(&BANDWIDTH, "bandwidth", BANDWIDTH,
		"bandwidth of prior belief")
	flag.IntVar(&TOTAL, "total", TOTAL,
		"total page count")
	flag.IntVar(&THIN, "thin", THIN,
		"beliefs are output once per 'thin' rows")
	flag.StringVar(&FLOATFMT, "floatFmt", FLOATFMT,
		"format for floats in the output CSV file")
}

// Creates the output CSV header
func makeHeader(m *model.Model) []string {
	total := len(m.Beliefs)
	header := make([]string, 3+2*total)
	header[0] = "iline"
	header[1] = "mean"
	header[2] = "variance"
	for i := 0; i != total; i++ {
		ipage := i + 1 // page numbering starts from 1
		header[3+2*i] = fmt.Sprintf("a%d", ipage)
		header[3+2*i+1] = fmt.Sprintf("b%d", ipage)
	}
	return header
}

// Creates a record from the model state
func makeRecord(iline int, m *model.Model) []string {
	total := len(m.Beliefs)
	mean, std := m.Avg()
	record := make([]string, 3+2*total)
	record[0] = fmt.Sprintf("%d", iline)
	record[1] = fmt.Sprintf(FLOATFMT, mean)
	record[2] = fmt.Sprintf(FLOATFMT, std)
	for i, b := range m.Beliefs {
		for j := 0; j != 2; j++ {
			record[3+2*i+j] = fmt.Sprintf(FLOATFMT, b[j])
		}
	}
	return record
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Create and initialize the model
	m := model.NewModel(TOTAL)
	m.Prior()

	// Go through the CSV data
	rdr := csv.NewReader(os.Stdin)
	wtr := csv.NewWriter(os.Stdout)
	defer wtr.Flush()

	// assume pps is the last column
	rdr.Read() // skip the header

	// write the output header
	header := makeHeader(m)
	err := wtr.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	// run through sessions and output beliefs every
	// 'thin' sessions
	for iline := 1; ; iline++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		pps, err := strconv.Atoi(record[len(record)-1])
		if err != nil {
			log.Printf("line %d: illegal pps %v, skipping",
				iline, record)
			continue
		}

		m.Update(BANDWIDTH, pps)
		if iline%THIN == 0 {
			record := makeRecord(iline, m)
			err := wtr.Write(record)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
