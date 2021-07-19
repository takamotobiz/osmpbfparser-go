package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	// ファイルを書き込み用にオープン (mode=0666)
	file, err := os.Create("./output.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	str := "Start:"
	str += time.Now().Format("2006-01-02 15:04:05")
	str += "\n"
	file.WriteString(str)
	fmt.Println(str)

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
				_, err = file.Write([]byte(rawJSON))
				if err != nil {
					log.Fatal(err)
				}
				_, err = file.Write([]byte("\n"))
				if err != nil {
					log.Fatal(err)
				}
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

	str1 := "End:"
	str1 += time.Now().Format("2006-01-02 15:04:05")
	str1 += "\n"
	file.WriteString(str1)
	fmt.Println(str1)

	if err := parser.Err(); err != nil {
		log.Fatal(err)
	}
}
