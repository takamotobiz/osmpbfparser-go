package bitmask

import (
	"encoding/gob"
	"fmt"
	"github.com/jneo8/logger-go"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
)

// PBFMasks - struct to hold common masks .
type PBFMasks struct {
	Nodes       *Bitmask
	Ways        *Bitmask
	Relations   *Bitmask
	WayRefs     *Bitmask
	RelNodes    *Bitmask
	RelWays     *Bitmask
	RelRelation *Bitmask
	Logger      *log.Logger
}

// NewPBFMasks - constructor
func NewPBFMasks() *PBFMasks {
	return &PBFMasks{
		Nodes:       NewBitMask(),
		Ways:        NewBitMask(),
		Relations:   NewBitMask(),
		WayRefs:     NewBitMask(),
		RelNodes:    NewBitMask(),
		RelWays:     NewBitMask(),
		RelRelation: NewBitMask(),
		Logger:      logger.NewLogger(),
	}
}

// WriteTo - write to destination
func (m *PBFMasks) WriteTo(sink io.Writer) (int64, error) {
	encoder := gob.NewEncoder(sink)
	err := encoder.Encode(m)
	return 0, err
}

// ReadFrom - read from destination
func (m *PBFMasks) ReadFrom(tap io.Reader) (int64, error) {
	decoder := gob.NewDecoder(tap)
	err := decoder.Decode(m)
	return 0, err
}

// WriteToFile - write to disk
func (m *PBFMasks) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if _, err := m.WriteTo(file); err != nil {
		return err
	}
	m.Logger.Debug("wrote bitmask:", path)
	return nil
}

// ReadFromFile - read from disk
func (m *PBFMasks) ReadFromFile(path string) error {

	// bitmask file doesn't exist
	if _, err := os.Stat(path); err != nil {
		fmt.Println("bitmask file not found:", path)
		os.Exit(1)
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	if _, err := m.ReadFrom(file); err != nil {
		return err
	}
	m.Logger.Debug("read bitmask:", path)
	return nil
}

// Print -- print debug stats
func (m PBFMasks) Print() {
	k := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	for i := 0; i < k.NumField(); i++ {
		key := k.Field(i).Name
		val := v.Field(i).Interface()
		fmt.Printf("%s: %v\n", key, (val.(*Bitmask)).Len())
	}
}
