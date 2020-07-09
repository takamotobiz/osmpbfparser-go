package osmpbfparser

import (
	"reflect"
	"testing"

	geojson "github.com/paulmach/go.geojson"
	"github.com/thomersch/gosmparse"
)

func TestBytesToElement(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Element
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesToElement(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesToElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesToElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToBytes(t *testing.T) {
	type fields struct {
		Type     int
		Node     gosmparse.Node
		Way      gosmparse.Way
		Relation gosmparse.Relation
		Elements []Element
		Role     int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Element{
				Type:     tt.fields.Type,
				Node:     tt.fields.Node,
				Way:      tt.fields.Way,
				Relation: tt.fields.Relation,
				Elements: tt.fields.Elements,
				Role:     tt.fields.Role,
			}
			got, err := e.ToBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToGeoJSON(t *testing.T) {
	type fields struct {
		Type     int
		Node     gosmparse.Node
		Way      gosmparse.Way
		Relation gosmparse.Relation
		Elements []Element
		Role     int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "basic node",
			fields: fields{
				Type: 0,
				Node: gosmparse.Node{
					Lat: 1.1,
					Lon: 2.2,
					Element: gosmparse.Element{
						ID: 10,
						Tags: map[string]string{
							"tagA": "A",
							"tagB": "B",
						},
					},
				},
			},
			want: []byte(`{"type":"FeatureCollection","features":[{"id":"node/10","type":"Feature","geometry":{"type":"Point","coordinates":[2.2,1.1]},"properties":{"osmType":"node","osmid":"node/10","tagA":"A","tagB":"B"}}]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Element{
				Type:     tt.fields.Type,
				Node:     tt.fields.Node,
				Way:      tt.fields.Way,
				Relation: tt.fields.Relation,
				Elements: tt.fields.Elements,
				Role:     tt.fields.Role,
			}
			if got := e.ToGeoJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.ToGeoJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestElement_nodeToJSON(t *testing.T) {
	type fields struct {
		Type     int
		Node     gosmparse.Node
		Way      gosmparse.Way
		Relation gosmparse.Relation
		Elements []Element
		Role     int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
		{
			name: "basic",
			fields: fields{
				Type: 0,
				Node: gosmparse.Node{
					Lat: 1.1,
					Lon: 2.2,
					Element: gosmparse.Element{
						ID: 10,
						Tags: map[string]string{
							"tagA": "A",
							"tagB": "B",
						},
					},
				},
			},
			want: []byte(`{"type":"FeatureCollection","features":[{"id":"node/10","type":"Feature","geometry":{"type":"Point","coordinates":[2.2,1.1]},"properties":{"osmType":"node","osmid":"node/10","tagA":"A","tagB":"B"}}]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Element{
				Type:     tt.fields.Type,
				Node:     tt.fields.Node,
				Way:      tt.fields.Way,
				Relation: tt.fields.Relation,
				Elements: tt.fields.Elements,
				Role:     tt.fields.Role,
			}
			if got := e.nodeToJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.nodeToJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestElement_NodeToFeature(t *testing.T) {
	type fields struct {
		Type     int
		Node     gosmparse.Node
		Way      gosmparse.Way
		Relation gosmparse.Relation
		Elements []Element
		Role     int
	}
	tests := []struct {
		name   string
		fields fields
		want   *geojson.Feature
	}{
		// TODO: Add test cases.
		{
			name: "basic",
			fields: fields{
				Type: 0,
				Node: gosmparse.Node{
					Lat: 1.1,
					Lon: 2.2,
					Element: gosmparse.Element{
						ID: 10,
						Tags: map[string]string{
							"tagA": "A",
							"tagB": "B",
						},
					},
				},
			},
			want: func() *geojson.Feature {

				f := geojson.NewPointFeature([]float64{2.2, 1.1})

				nodeID := "node/10"
				f.ID = nodeID
				f.SetProperty("osmid", "node/10")
				f.SetProperty("osmType", "node")

				for k, v := range map[string]string{
					"tagA": "A",
					"tagB": "B",
				} {
					f.SetProperty(k, v)
				}
				return f
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Element{
				Type:     tt.fields.Type,
				Node:     tt.fields.Node,
				Way:      tt.fields.Way,
				Relation: tt.fields.Relation,
				Elements: tt.fields.Elements,
				Role:     tt.fields.Role,
			}
			if got := e.NodeToFeature(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.NodeToFeature() = %v, want %v", got, tt.want)
			}
		})
	}
}
