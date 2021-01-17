package progressbar

import (
	"fmt"
	"io"

	download "github.com/admpub/go-download/v2"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

func New(opt *download.Options, width ...int) *mpb.Progress {
	var w int
	if len(width) > 0 {
		w = width[0]
	}
	if w <= 0 {
		w = 80
	}
	progress := mpb.New(mpb.WithWidth(w))
	//defer progress.Wait()
	opt.Proxy = func(name string, download int, size int64, r io.Reader) io.Reader {
		name = fmt.Sprintf("%s-%d", name, download)
		bar := progress.AddBar(
			size,
			mpb.PrependDecorators(
				decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
				decor.CountersNoUnit(`%3d / %3d`, decor.WC{W: 18}),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WC{W: 5}),
			),
		)
		return bar.ProxyReader(r)
	}
	return progress
}
