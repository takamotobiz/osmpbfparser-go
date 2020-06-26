package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type pbfParser struct {
	PBFMasks *bitmask.PBFMasks

	// leveldb
	LevelDB *leveldb.DB
	Args    Args

	// Log
	Logger *log.Logger
}

// Run ...
func (p *pbfParser) Run() error {
	p.Logger.Infof("%+v", p)
	db, err := leveldb.OpenFile(
		p.Args.LevelDBPath,
		&opt.Options{DisableBlockCache: true},
	)
	if err != nil {
		p.Logger.Error(err)
		return err
	}
	defer db.Close()

	indexer := newPBFIndexer(p.Args.PBFFile, p.PBFMasks)
	if err := indexer.Run(); err != nil {
		return err
	}

	return nil
}

// SetLogger ...
func (p *pbfParser) SetLogger(logger *log.Logger) {
	p.Logger = logger
}
