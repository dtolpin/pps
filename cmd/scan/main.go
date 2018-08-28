// Command scan runs over pps sequence (stored as a CSV file, last called is the pps)
// and computes updated beliefs after each session. The output is enumerated CSV, where
// each line is a flattened Beliefs vector.
package main

import (
	csv "encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"bitbucket.org/dtolpin/pps/model"
)

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
func makeRecord(iline int, m *model.Model, floatFmt string) []string {
	total := len(m.Beliefs)
	mean, std := m.Avg()
	record := make([]string, 3 + 2*total)
	record[0] = fmt.Sprintf("%d", iline)
	record[1] = fmt.Sprintf(floatFmt, mean)
	record[2] = fmt.Sprintf(floatFmt, std)
	for i, b := range m.Beliefs {
		for j := 0; j != 2; j++ {
			record[3+2*i+j] = fmt.Sprintf(floatFmt, b[j])
		}
	}
	return record
}

func main() {
	bandwidth := *flag.Float64("bandwidth", 100.,
		"bandwidth of prior belief")
	total := *flag.Int("total", 10,
		"total page count")
	thin := *flag.Int("thin", 100,
		"beliefs are output once per 'thin' rows")
	floatFmt := *flag.String("floatFmt", "%.3f",
		"format for floats in the output CSV file")
	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Create and initialize the model
	m := model.NewModel(total)
	m.Prior()

	// Go through the CSV data
	rdr := csv.NewReader(os.Stdin)
	wtr := csv.NewWriter(os.Stdout)
	wtr.Flush()

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

		m.Update(bandwidth, pps)
		if iline % thin == 0 {
			record := makeRecord(iline, m, floatFmt)
			err := wtr.Write(record)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
