package osmpbfparser

import (
	"github.com/paulmach/go.geojson"
	"strconv"
)

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
