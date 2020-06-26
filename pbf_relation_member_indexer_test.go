package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	"testing"
)

func TestPBFRelationMemberIndexer(t *testing.T) {
	masks := bitmask.NewPBFMasks()
	indexer := newPBFIndexer(
		"./assert/test.pbf",
		masks,
	)
	if err := indexer.Run(); err != nil {
		t.Error(err)
	}

	rmIndexer := newPBFRelationMemberIndexer(
		"./assert/test.pbf",
		masks,
	)
	if err := rmIndexer.Run(); err != nil {
		t.Error(err)
	}
}
