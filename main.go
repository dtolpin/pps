package main

import (
	"flag"
	"log"

	"bitbucket.org/dtolpin/pps/infr"
)

const maxargs = 1

func main() {
	bandwidth := flag.Float64("bandwidth", 100.,
		"bandwidth of prior belief")
	total := flag.Int("total", 10,
		"total page count")
	flag.Parse()
	log.Printf("bandwidth=%v total=%v", *bandwidth, *total)

	args := flag.Args()
	if len(args) > 1 {
		log.Fatalf("too many args, expected at most %v, got %v: %v",
			maxargs, len(args), args)
	}

	var m infr.Model
	m.Init(*total)
}
