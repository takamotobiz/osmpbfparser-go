# OSM PBF Parser

- Parse pbf file to osm element(Node, Way and Relation)

- Iterator pattern. Easy to get.

- Progress bar for parsing process.

- Control memory usage by using levelDB.

- Convert Element to GeoJSON
    - Support type: Node, Way, Relation(polygon, multipolygon)

## Install

```
go get -u https://github.com/jneo8/osmpbfparser-go
```

## Quick start

- **Download the pbf file**

```go
package main

import (
	"fmt"
	"github.com/jneo8/osmpbfparser-go"
	"log"
)

func main() {
	parser := osmpbfparser.New(
		osmpbfparser.Args{
			PBFFile:     "./assert/test.pbf",
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
```
