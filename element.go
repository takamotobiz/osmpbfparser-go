package osmpbfparser

import (
	"bytes"
	"encoding/gob"
	"github.com/paulmach/go.geojson"
	"github.com/thomersch/gosmparse"
	"strconv"
)

// BytesToElement convert bytes to Element struct.
func BytesToElement(b []byte) (Element, error) {
	decoder := gob.NewDecoder(bytes.NewReader(b))
	var element Element
	err := decoder.Decode(&element)
	return element, err
}

// Element is an osm data element group set.
type Element struct {
	Type     int // 0=Node, 1=Way. 2=Relation
	Node     gosmparse.Node
	Way      gosmparse.Way
	Relation gosmparse.Relation
	Elements []Element
	Role     int // 0=outer, 1=inner
}

// ToBytes convert element struct to bytes.
func (e *Element) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(e); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// ToGeoJSON convery element to JSON bytes.
func (e *Element) ToGeoJSON() []byte {
	var b []byte
	switch e.Type {
	case 0:
		b = e.nodeToJSON()
	}
	return b
}

func (e *Element) nodeToJSON() []byte {
	fc := geojson.NewFeatureCollection()
	fc.AddFeature(e.NodeToFeature())
	rawJSON, _ := fc.MarshalJSON()
	return rawJSON
}

// NodeToFeature convert node element to geojson feature.
func (e *Element) NodeToFeature() *geojson.Feature {
	f := geojson.NewPointFeature([]float64{e.Node.Lon, e.Node.Lat})

	nodeID := "node/" + strconv.FormatInt(e.Node.ID, 10)
	f.ID = nodeID
	f.SetProperty("osmid", nodeID)
	f.SetProperty("osmType", "node")

	for k, v := range e.Node.Tags {
		f.SetProperty(k, v)
	}
	return f
}
