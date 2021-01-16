package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	download "github.com/admpub/go-download/v2"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"github.com/webx-top/com/httpClientOptions"
)

// NewJar Cookie record Jar
func newJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(nil)
	return cookieJar
}

func main() {

	url := os.Args[len(os.Args)-1]

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
		Client: func() http.Client {
			client := download.NewHTTPClient(
				httpClientOptions.CheckRedirect(func(req *http.Request, via []*http.Request) error {
					log.Printf("Redirect: %v\n", req.URL)
					return nil
				}),
				httpClientOptions.CookieJar(newJar()),
				httpClientOptions.InsecureSkipVerify(),
			)
			return *client
		},
	}

	f, err := download.Open(url, options)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var output *os.File
	name := info.Name()
	output, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	if _, err = io.Copy(output, f); err != nil {
		log.Fatal(err)
	}

	log.Printf("Success. %s saved.\n", name)
}
