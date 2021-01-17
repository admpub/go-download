package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	download "github.com/admpub/go-download/v2"
	"github.com/admpub/go-download/v2/progressbar"
	"github.com/webx-top/com/httpClientOptions"
)

// NewJar Cookie record Jar
func newJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(nil)
	return cookieJar
}

func main() {
	url := os.Args[len(os.Args)-1]
	options := &download.Options{
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
	progress := progressbar.New(options)
	defer progress.Wait()
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
