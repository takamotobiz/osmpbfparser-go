package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type pbfParser struct {
	PBFMasks *bitmask.PBFMasks

	// leveldb
	LevelDB *leveldb.DB
	Args    Args
}

// Run ...
func (p *pbfParser) Run() error {
	db, err := leveldb.OpenFile(
		p.Args.LevelDBPath,
		&opt.Options{DisableBlockCache: true},
	)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer db.Close()
	return nil
}
