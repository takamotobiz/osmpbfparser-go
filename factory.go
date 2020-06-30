package osmpbfparser

import (
	"github.com/jneo8/logger-go"
	"github.com/jneo8/osmpbfparser-go/bitmask"
)

// New ...
func New(
	args Args,
) PBFParser {
	return &pbfParser{
		Args:   args,
		Logger: logger.NewLogger(),
	}
}

func newPBFIndexer(
	pbfFile string,
	pbfMasks *bitmask.PBFMasks,
) pbfDataParser {
	return &PBFIndexer{
		PBFFile:  pbfFile,
		PBFMasks: pbfMasks,
	}
}

func newPBFRelationMemberIndexer(
	pbfFile string,
	pbfMasks *bitmask.PBFMasks,
) pbfDataParser {
	return &PBFRelationMemberIndexer{
		PBFFile:  pbfFile,
		PBFMasks: pbfMasks,
	}
}
