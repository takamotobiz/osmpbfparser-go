package main

import (
	"fmt"
	"log"

	"github.com/jneo8/osmpbfparser-go"
)

func main() {
	parser := osmpbfparser.New(
		osmpbfparser.Args{
			PBFFile:     "/Users/takamotokeiji/Downloads/shikoku-latest.osm.pbf",
			LevelDBPath: "/tmp/osmpbfparser",
			BatchSize:   10000,
		},
	)
	var nc, wc, rc uint64

	for emt := range parser.Iterator() {
		// rawJSON, err := emt.ToGeoJSON()
		// if err != nil {
		//     log.Fatal(err)
		// }
		// fmt.Println(string(rawJSON))

		tags := emt.GetTags()
		if nv, fl := tags["amenity"]; fl == true {
			if nv == "school" {
				rawJSON, err := emt.ToGeoJSON()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(rawJSON))
				switch emt.Type {
				case 0: //Node
					nc++
				case 1: //Way
					wc++
				case 2: //Relation
					rc++
				}
			}
		}
	}
	fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)

	if err := parser.Err(); err != nil {
		log.Fatal(err)
	}
}
