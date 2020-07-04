package osmpbfparser

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
)

func newBar(n int, prefix string) *pb.ProgressBar {
	bar := pb.StartNew(n)
	bar.SetWidth(100)
	bar.SetMaxWidth(-1)
	bar.SetTemplate(pb.Full)
	bar.Set("suffix", "   \n")
	bar.Set("prefix", fmt.Sprintf("%s: ", prefix))
	return bar
}
