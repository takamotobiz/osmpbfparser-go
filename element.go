package osmpbfparser

import (
	"bytes"
	"encoding/gob"
	"github.com/thomersch/gosmparse"
)

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

// BytesToElement convert bytes to Element struct.
func BytesToElement(b []byte) (Element, error) {
	decoder := gob.NewDecoder(bytes.NewReader(b))
	var element Element
	err := decoder.Decode(&element)
	return element, err
}
