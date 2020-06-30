package osmpbfparser

import (
	"bytes"
	"encoding/binary"
	"github.com/jneo8/osmpbfparser-go/bitmask"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/thomersch/gosmparse"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
)

type pbfParser struct {
	PBFMasks *bitmask.PBFMasks

	// leveldb
	LevelDB *leveldb.DB
	Args    Args

	// Log
	Logger      *log.Logger
	elementChan chan Element

	Error error
}

func (p *pbfParser) Err() error {
	return p.Error
}

// Run ...
func (p *pbfParser) Iterator() <-chan Element {
	p.Logger.Infof("%+v", p)

	outputCh := make(chan Element)

	go func() {
		defer close(outputCh)
		db, err := leveldb.OpenFile(
			p.Args.LevelDBPath,
			&opt.Options{DisableBlockCache: true},
		)
		if err != nil {
			p.Logger.Error(err)
			p.Error = err
			return
		}
		defer db.Close()
		p.LevelDB = db
		p.Logger.Info("Init leveldb")

		// bitmask
		p.PBFMasks = bitmask.NewPBFMasks()

		// Index
		indexer := newPBFIndexer(p.Args.PBFFile, p.PBFMasks)
		if err := indexer.Run(); err != nil {
			p.Error = err
			return
		}
		// Relation member indexer
		relationMemberIndexer := newPBFRelationMemberIndexer(p.Args.PBFFile, p.PBFMasks)
		if err := relationMemberIndexer.Run(); err != nil {
			p.Error = err
			return
		}

		p.Logger.Info("Finish index")

		reader, err := os.Open(p.Args.PBFFile)
		if err != nil {
			p.Error = err
			return
		}
		defer reader.Close()

		// FirstRound
		// Put way refs, relation member into db.
		batch := leveldb.MakeBatch(p.Args.FlushSize * 1024 * 1024)
		p.elementChan = make(chan Element)

		firstRoundWg := sync.WaitGroup{}
		firstRoundWg.Add(1)
		errCount := make(map[int]int)
		go func() {
			defer firstRoundWg.Done()
			for emt := range p.elementChan {
				switch emt.Type {
				case 0:
					if p.PBFMasks.WayRefs.Has(emt.Node.ID) || p.PBFMasks.RelNodes.Has(emt.Node.ID) {
						id, b := nodeToBytes(emt.Node)
						batch.Put(
							[]byte(id),
							b,
						)
					}
				case 1:
					if p.PBFMasks.Ways.Has(emt.Way.ID) {
						emtBytes, err := emt.ToBytes()
						if err != nil {
							errCount[1]++
							continue
						}
						batch.Put(
							[]byte("W"+strconv.FormatInt(emt.Way.ID, 10)),
							emtBytes,
						)
					}
				case 2:
					if p.PBFMasks.RelRelation.Has(emt.Relation.ID) {
						emtBytes, err := emt.ToBytes()
						if err != nil {
							errCount[2]++
							continue
						}
						batch.Put(
							[]byte("R"+strconv.FormatInt(emt.Relation.ID, 10)),
							emtBytes,
						)

					}
				}
			}
		}()
		firstRoundDecoder := gosmparse.NewDecoder(reader)
		if err := firstRoundDecoder.Parse(p); err != nil {
			p.Error = err
			return
		}
		close(p.elementChan)
		firstRoundWg.Wait()
		p.Logger.Info("Finish first round")
		if _, err := reader.Seek(0, io.SeekStart); err != nil {
			p.Error = err
			return
		}

		// Final round.
		// Re-init element chan
		p.elementChan = make(chan Element)

		finalRoundWg := sync.WaitGroup{}
		finalRoundWg.Add(1)
		go func() {
			defer finalRoundWg.Done()
			for emt := range p.elementChan {
				outputCh <- emt
			}
		}()

		decoder := gosmparse.NewDecoder(reader)
		if err := decoder.Parse(p); err != nil {
			p.Error = err
			return
		}
		close(p.elementChan)
		finalRoundWg.Wait()
	}()

	return outputCh
}

// ReadNode ...
func (p *pbfParser) ReadNode(node gosmparse.Node) {
	p.elementChan <- Element{
		Type: 0,
		Node: node,
	}
}

// ReadWay ...
func (p *pbfParser) ReadWay(way gosmparse.Way) {
	p.elementChan <- Element{
		Type: 1,
		Way:  way,
	}
}

// ReadRelation ...
func (p *pbfParser) ReadRelation(relation gosmparse.Relation) {
	p.elementChan <- Element{
		Type:     2,
		Relation: relation,
	}
}

// SetLogger ...
func (p *pbfParser) SetLogger(logger *log.Logger) {
	p.Logger = logger
}

func nodeToBytes(n gosmparse.Node) (string, []byte) {
	var buf bytes.Buffer

	var latBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(latBytes, math.Float64bits(n.Lat))
	buf.Write(latBytes)
	return strconv.FormatInt(n.ID, 10), buf.Bytes()
}
