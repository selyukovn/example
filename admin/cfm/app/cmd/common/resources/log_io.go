package resources

import (
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"os"
)

func NewLogIoFile(pathToFile string) io.WriteCloser {
	assert.Str().NotEmpty().Must(pathToFile)

	file, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return file
}
