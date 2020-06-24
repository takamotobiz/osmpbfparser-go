package osmpbfparser

import (
	"testing"

	"github.com/jneo8/osmpbfparser-go/bitmask"
	"github.com/syndtr/goleveldb/leveldb"
)

func Test_pbfParser_Run(t *testing.T) {
	type fields struct {
		PBFMasks *bitmask.PBFMasks
		LevelDB  *leveldb.DB
		Args     Args
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
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pbfParser{
				PBFMasks: tt.fields.PBFMasks,
				LevelDB:  tt.fields.LevelDB,
				Args:     tt.fields.Args,
			}
			if err := p.Run(); (err != nil) != tt.wantErr {
				t.Errorf("pbfParser.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
