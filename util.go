package download

import (
	"io"
	"os"
)

func Download(url, saveName string, options *Options) (int64, error) {
	f, err := Open(url, options)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	var output *os.File
	if len(saveName) == 0 {
		saveName = info.Name()
	}
	output, err = os.OpenFile(saveName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	return io.Copy(output, f)
}
