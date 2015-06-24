package reader

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/jksy/go-log-viewer/reader"
)

type withTempfileFunc func(f *os.File)

func TestFileAppend(t *testing.T) {
	withTempfile(t, func(testFile *os.File) {
		ch := make(chan []byte)
		fileReader := reader.NewFileReaderWithOpen(testFile.Name(), ch)
		go fileReader.Run()

		str := randomString()
		bin := []byte(str)
		testFile.Write(bin)

		select {
		case readed := <-ch:
			if string(readed[:]) != str {
				t.Error("dont same string, %v:%v", str, string(readed[:]))
			}
		case <-time.After(1 * time.Second):
			t.Error("cant read channel from file")
		}
	})
}

func withTempfile(t *testing.T, f withTempfileFunc) {
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Error("cant create tempfile")
	}
	f(file)
	defer os.Remove(file.Name())
}

func randomString() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d", time.Now()))
	return fmt.Sprintf("%x", h.Sum(nil))
}
