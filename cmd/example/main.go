package main

import (
	"fmt"
	"github.com/jneo8/osmpbfparser-go"
	"log"
)

func main() {
	parser := osmpbfparser.New(
		osmpbfparser.Args{
			PBFFile:     "./static/taiwan-latest.osm.pbf",
			LevelDBPath: "/tmp/osmpbfparser",
			BatchSize:   10000,
		},
	)

	for emt := range parser.Iterator() {
		rawJSON, err := emt.ToGeoJSON()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(rawJSON))
	}
	if err := parser.Err(); err != nil {
		log.Fatal(err)
	}
}
