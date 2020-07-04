package main

import (
	"github.com/jneo8/osmpbfparser-go"
	"log"
)

func main() {
	parser := osmpbfparser.New(
		osmpbfparser.Args{
			PBFFile:     "./assert/taiwan-latest.osm.pbf",
			LevelDBPath: "/tmp/osmpbfparser",
			BatchSize:   10000,
		},
	)

	for range parser.Iterator() {
		// fmt.Println(emt)
	}
	if err := parser.Err(); err != nil {
		log.Fatal(err)
	}
}
