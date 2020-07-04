package osmpbfparser

import (
	"github.com/cheggaaa/pb/v3"
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
		Report: Report{PBFFile: args.PBFFile},
	}
}

func newPBFIndexer(
	pbfFile string,
	pbfMasks *bitmask.PBFMasks,
	bar *pb.ProgressBar,
) pbfDataParser {
	return &PBFIndexer{
		PBFFile:  pbfFile,
		PBFMasks: pbfMasks,
		Bar:      bar,
	}
}

func newPBFRelationMemberIndexer(
	pbfFile string,
	pbfMasks *bitmask.PBFMasks,
	bar *pb.ProgressBar,
) pbfDataParser {
	return &PBFRelationMemberIndexer{
		PBFFile:  pbfFile,
		PBFMasks: pbfMasks,
		Bar:      bar,
	}
}

func newPBFCounter(
	pbfFile string,
) pbfDataCounter {
	return &PBFCounter{
		PBFFile: pbfFile,
	}
}
