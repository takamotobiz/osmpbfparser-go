package osmpbfparser

import (
	log "github.com/sirupsen/logrus"
	"github.com/thomersch/gosmparse"
)

// PBFParser ...
type PBFParser interface {
	Iterator() <-chan Element
	SetLogger(*log.Logger)
	Err() error
}

type pbfDataParser interface {
	gosmparse.OSMReader
	Run() error
}

// Reader ...
type Reader interface {
	Run(emt Element)
}
