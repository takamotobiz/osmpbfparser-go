package osmpbfparser

import (
	log "github.com/sirupsen/logrus"
	"github.com/thomersch/gosmparse"
)

// PBFParser ...
type PBFParser interface {
	Run() error
	SetLogger(*log.Logger)
}

type pbfDataParser interface {
	gosmparse.OSMReader
	Run() error
}
