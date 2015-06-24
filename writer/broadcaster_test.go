package writer

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/jksy/go-log-viewer/writer"
)

type TestWriter struct {
	Writed []byte
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	w.Writed = make([]byte, len(p))
	copy(w.Writed[:], p[:])
	return len(p), nil
}

func TestWrite(t *testing.T) {
	ch := make(chan []byte)
	broadcaster := writer.NewBroadcaster(ch)
	testWriter := TestWriter{}
	broadcaster.Add(&testWriter)

	str := randomString()
	bin := []byte(str)
	go broadcaster.Run()
	broadcaster.Write(bin)
	defer func() {
		broadcaster.Close()
	}()
	time.Sleep(20 * time.Millisecond) // wait for writing data

	if str != string(testWriter.Writed[:]) {
		t.Errorf("dont same string, %v, %v", str, string(testWriter.Writed[:]))
	}
}

func randomString() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d", time.Now()))
	return fmt.Sprintf("%x", h.Sum(nil))
}
