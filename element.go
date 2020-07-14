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

// IsArea return if element is area or not.
// https://wiki.openstreetmap.org/wiki/Key:area
func (e *Element) IsArea() bool {
	var isPolygon bool
	if val, ok := e.Way.Tags["area"]; ok && val == "yes" {
		// This list is probably incomplete - please add other cases
		if _, ok := e.Way.Tags["highway"]; ok {
			// highway=*, see 'Highway areas' below for more details
		} else if _, isBarrier := e.Way.Tags["barrier"]; isBarrier {
			// barrier=*, for thicker hedges or walls or detailed mapping defined using an area add area=yes
		} else {
			isPolygon = true
		}
	}
	return isPolygon
}

// ToGeoJSON convery element to JSON bytes.
func (e *Element) ToGeoJSON() []byte {
	var b []byte
	switch e.Type {
	case 0:
		b = e.nodeToJSON()
	case 1:
		b = e.wayToJSON()
	}
	return b
}

func (e *Element) nodeToJSON() []byte {
	fc := geojson.NewFeatureCollection()
	fc.AddFeature(e.NodeToFeature())
	rawJSON, _ := fc.MarshalJSON()
	return rawJSON
}

func (e *Element) wayToJSON() []byte {
	fc := geojson.NewFeatureCollection()
	fc.AddFeature(e.NodeToFeature())
	rawJSON, _ := fc.MarshalJSON()
	return rawJSON
}

// WayToJSON convert way element to geojson feature.
func (e *Element) WayToJSON() *geojson.Feature {
	latLngs := [][]float64{}
	for _, member := range e.Elements {
		latLngs = append(latLngs, []float64{member.Node.Lon, member.Node.Lat})
	}

	var f *geojson.Feature

	switch e.IsArea() {
	case true:
		f = geojson.NewPolygonFeature([][][]float64{latLngs})
	default:
		f = geojson.NewLineStringFeature(latLngs)
	}

	wayID := "way" + "/" + strconv.FormatInt(e.Way.ID, 10)
	f.ID = wayID
	f.SetProperty("osmid", wayID)
	f.SetProperty("osmType", "way")

	for k, v := range e.Way.Tags {
		f.SetProperty(k, v)
	}
	return f
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
