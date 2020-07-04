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
		newBar(1166412, "Indexer"),
	)
	if err := indexer.Run(); err != nil {
		t.Error(err)
	}

	rmIndexer := newPBFRelationMemberIndexer(
		"./assert/test.pbf",
		masks,
		newBar(1166412, "RM Indexer"),
	)
	if err := rmIndexer.Run(); err != nil {
		t.Error(err)
	}
}
