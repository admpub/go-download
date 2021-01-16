package main

import (
	"fmt"
	"io"
	"log"

	download "github.com/admpub/go-download/v2"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

func main() {

	progress := mpb.New(mpb.WithWidth(80))
	defer progress.Wait()

	options := &download.Options{
		Proxy: func(name string, download int, size int64, r io.Reader) io.Reader {
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
		},
	}

	f, err := download.Open("https://storage.googleapis.com/golang/go1.8.1.src.tar.gz", options)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// f implements io.Reader, write file somewhere or do some other sort of work with it
}
