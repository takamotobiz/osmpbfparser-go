package osmpbfparser

import (
	"github.com/jneo8/osmpbfparser-go/bitmask"
	"github.com/thomersch/gosmparse"
	"github.com/vbauerster/mpb/v5"
	"os"
)

// PBFRelationMemberIndexer ...
type PBFRelationMemberIndexer struct {
	PBFFile     string
	PBFMasks    *bitmask.PBFMasks
	NodeBar     *mpb.Bar
	WayBar      *mpb.Bar
	RelationBar *mpb.Bar
}

// Run ...
func (p *PBFRelationMemberIndexer) Run() error {
	reader, err := os.Open(p.PBFFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := gosmparse.NewDecoder(reader)
	if err := decoder.Parse(p); err != nil {
		return err
	}
	return nil
}

// ReadNode ...
func (p *PBFRelationMemberIndexer) ReadNode(node gosmparse.Node) {
	defer p.NodeBar.Increment()
}

// ReadWay ...
func (p *PBFRelationMemberIndexer) ReadWay(way gosmparse.Way) {
	defer p.WayBar.Increment()
	if p.PBFMasks.RelWays.Has(way.ID) {
		for _, nodeID := range way.NodeIDs {
			p.PBFMasks.RelNodes.Insert(nodeID)
		}
	}
}

// ReadRelation ...
func (p *PBFRelationMemberIndexer) ReadRelation(relation gosmparse.Relation) {
	defer p.RelationBar.Increment()
}
