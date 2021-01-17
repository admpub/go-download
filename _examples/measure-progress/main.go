package main

import (
	"log"

	download "github.com/admpub/go-download/v2"
	"github.com/admpub/go-download/v2/progressbar"
)

func main() {
	options := &download.Options{}
	progress := progressbar.New(options)
	defer progress.Wait()
	f, err := download.Open("https://storage.googleapis.com/golang/go1.8.1.src.tar.gz", options)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// f implements io.Reader, write file somewhere or do some other sort of work with it
}
