package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/thomersch/gosmparse"
	"io"
	"os"
	"strconv"
	"sync"
)

type pbfParser struct {
	PBFMasks *bitmask.PBFMasks

	// leveldb
	DB    *leveldb.DB
	Batch *leveldb.Batch
	Args  Args

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
		p.DB = db
		p.Logger.Infof("Init leveldb at %s", p.Args.LevelDBPath)

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
		p.Batch = leveldb.MakeBatch(p.Args.BatchSize)
		defer p.Batch.Reset()

		p.elementChan = make(chan Element)

		firstRoundWg := sync.WaitGroup{}
		firstRoundWg.Add(1)
		errCount := make(map[int]int)
		go func() {
			defer firstRoundWg.Done()
			idx := 0
			for emt := range p.elementChan {
				idx++
				// Check batch size every batchsize.
				if idx%p.Args.BatchSize == 0 {
					if err := p.checkBatch(); err != nil {
						p.Error = err
					}
				}
				switch emt.Type {
				case 0:
					if p.PBFMasks.WayRefs.Has(emt.Node.ID) || p.PBFMasks.RelNodes.Has(emt.Node.ID) {
						id, b := nodeToBytes(emt.Node)
						p.Batch.Put(
							[]byte(id),
							b,
						)
					}
				case 1:
					if p.PBFMasks.RelWays.Has(emt.Way.ID) {
						emtBytes, err := emt.ToBytes()
						if err != nil {
							errCount[1]++
							continue
						}
						p.Batch.Put(
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
						p.Batch.Put(
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

		// Flush batch.
		if err := p.flushBatch(true); err != nil {
			p.Error = err
			return
		}
		p.Logger.Info("Finish first round")

		// Rewind pbf file reader
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
				switch emt.Type {
				case 0:
					if p.PBFMasks.Nodes.Has(emt.Node.ID) {
						outputCh <- emt
					}
				case 1:
					if p.PBFMasks.Ways.Has(emt.Way.ID) {
						emts, err := p.dbLookupWayElement(&emt.Way)
						if err != nil {
							p.Logger.Warning(err)
							continue
						}
						emt.Elements = emts
						outputCh <- emt
					}
				}
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

func (p *pbfParser) dbLookupWayElement(way *gosmparse.Way) ([]Element, error) {
	var emts []Element
	for _, nodeID := range way.NodeIDs {
		id := strconv.FormatInt(nodeID, 10)
		b, err := p.DB.Get(
			[]byte(id),
			nil,
		)
		if err != nil {
			return emts, err
		}
		node := bytesToNode(b)
		emt := Element{
			Type: 0,
			Node: node,
		}
		emts = append(emts, emt)
	}
	return emts, nil
}

func (p *pbfParser) checkBatch() error {
	if p.Batch.Len() >= p.Args.BatchSize {
		if err := p.flushBatch(true); err != nil {
			return err
		}
	}
	return nil
}

func (p *pbfParser) flushBatch(sync bool) error {
	writeOpts := opt.WriteOptions{
		NoWriteMerge: true,
		Sync:         sync,
	}
	if err := p.DB.Write(p.Batch, &writeOpts); err != nil {
		return err
	}
	p.Batch.Reset()
	return nil
}
