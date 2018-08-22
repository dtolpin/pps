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

const maxargs = 1

func main() {
	bandwidth := flag.Float64("bandwidth", 100.,
		"bandwidth of prior belief")
	total := flag.Int("total", 10,
		"total page count")
	thin := flag.Int("thin", 100,
		"beliefs are output once per 'thin' rows")
	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatalf("unexpected position arguments: %v", flag.Args())
	}

	// Create and initialize the model
    m := model.NewModel(*total)
	m.Prior()

	// Go through the CSV data
	rdr := csv.NewReader(os.Stdin)
	wtr := csv.NewWriter(os.Stdout)

    rdr.Read() // skip the header
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

		m.Update(*bandwidth, pps)
		if iline % *thin == 0 {
			record := make([]string, 2*len(m.Beliefs) + 1)
			record[0] = strconv.Itoa(iline)
			for i, b := range m.Beliefs {
				for j := 0; j != 2; j++ {
					record[1 + 2*i + j] = fmt.Sprintf("%.2f", b[j])
				}
			}
            err := wtr.Write(record)
            if err != nil {
                log.Fatal(err)
            }
		}
	}

    wtr.Flush()
}
