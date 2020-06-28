package osmpbfparser

import (
	"github.com/jneo8/logger-go"
	"github.com/jneo8/osmpbfparser-go/bitmask"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func Test_pbfParser_Run(t *testing.T) {
	type fields struct {
		PBFMasks *bitmask.PBFMasks
		LevelDB  *leveldb.DB
		Args     Args
		Logger   *log.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				Args: Args{
					LevelDBPath: "/tmp/osmpbfparser",
					PBFFile:     "./assert/test.pbf",
				},
				Logger:   logger.NewLogger(),
				PBFMasks: bitmask.NewPBFMasks(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pbfParser{
				PBFMasks: tt.fields.PBFMasks,
				LevelDB:  tt.fields.LevelDB,
				Args:     tt.fields.Args,
				Logger:   tt.fields.Logger,
			}
			if err := p.Run(); (err != nil) != tt.wantErr {
				t.Errorf("pbfParser.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
