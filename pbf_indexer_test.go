package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	"testing"
)

func TestPBFIndexer(t *testing.T) {
	indexer := newPBFIndexer(
		"./assert/test.pbf",
		bitmask.NewPBFMasks(),
	)
	if err := indexer.Run(); err != nil {
		t.Error(err)
	}
}
