package osmpbfparser

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/jneo8/osmpbfparser-go/bitmask"
	"github.com/thomersch/gosmparse"
	"os"
)

// PBFIndexer ...
type PBFIndexer struct {
	PBFFile  string
	PBFMasks *bitmask.PBFMasks
	Bar      *pb.ProgressBar
}

// Run ...
func (p *PBFIndexer) Run() error {
	reader, err := os.Open(p.PBFFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	defer p.Bar.Finish()

	decoder := gosmparse.NewDecoder(reader)
	if err := decoder.Parse(p); err != nil {
		return err
	}
	return nil
}

// ReadNode ...
func (p *PBFIndexer) ReadNode(node gosmparse.Node) {
	defer p.Bar.Increment()
	// Get node if tags > 0
	if len(node.Tags) == 0 {
		return
	}
	p.PBFMasks.Nodes.Insert(node.ID)
}

// ReadWay ...
func (p *PBFIndexer) ReadWay(way gosmparse.Way) {
	defer p.Bar.Increment()
	if len(way.Tags) == 0 {
		return
	}
	p.PBFMasks.Ways.Insert(way.ID)
	for _, nodeID := range way.NodeIDs {
		p.PBFMasks.WayRefs.Insert(nodeID)
	}
}

// ReadRelation ...
func (p *PBFIndexer) ReadRelation(relation gosmparse.Relation) {
	defer p.Bar.Increment()
	if len(relation.Tags) == 0 {
		return
	}
	var count = make(map[int]int64)
	for _, member := range relation.Members {
		count[int(member.Type)]++
	}
	// Skip if relations cotain 0 way
	if count[1] == 0 {
		return
	}
	p.PBFMasks.Relations.Insert(relation.ID)
	for _, member := range relation.Members {
		switch member.Type {
		case 0:
			p.PBFMasks.RelNodes.Insert(member.ID)
		case 1:
			p.PBFMasks.RelWays.Insert(member.ID)
		case 2:
			p.PBFMasks.RelRelation.Insert(member.ID)
		}
	}
}
