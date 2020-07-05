package osmpbfparser

import (
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"sync"
)

func newProgress(n int) (*mpb.Progress, *sync.WaitGroup) {
	var wg sync.WaitGroup
	wg.Add(n)
	return mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100)), &wg
}

func addBar(pb *mpb.Progress, name string, total int) *mpb.Bar {
	return pb.AddBar(
		int64(total),
		mpb.BarStyle("╢▌▌░╟"),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
		),
		mpb.AppendDecorators(
			// decor.OnComplete(
			// 	decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 5, C: decor.DidentRight}),
			// 	"done",
			// ),
			// decor.Counters(1, "% d / % d", decor.WC{W: 5, C: decor.DidentRight}),
			decor.CountersNoUnit("%d / %d ", decor.WCSyncWidth),

			decor.Percentage(),
		),
	)
}
