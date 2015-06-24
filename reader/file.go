package reader

import (
	"github.com/go-fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
)

type File struct {
	watcher  *fsnotify.Watcher
	filename string
	fp       *os.File
	ch       chan []byte
	closed   chan bool
}

func NewFileReader() *File {
	file := &File{}
	file.closed = make(chan bool)
	return file
}

func NewFileReaderWithOpen(path string, ch chan []byte) *File {
	file := NewFileReader()
	file.open(path, ch)
	return file
}

func (f *File) open(spec string, ch chan []byte) {
	f.openFile(spec)
	f.watchFile(spec)
	f.ch = ch
	log.Println("watching... " + spec + "")
}

func (f *File) Reopen() {
	f.open(f.filename, f.ch)
}

func (f *File) Run() {
	for {
		select {
		case ev := <-f.watcher.Events:
			if ev.Op&fsnotify.Write != 0 {
				buf, err := ioutil.ReadAll(f.fp)
				if err != nil {
					f.Reopen()
				} else {
					f.ch <- buf
				}
			}
		case err := <-f.watcher.Errors:
			log.Println("error:", err)
			f.Close()
			return
		case <-f.closed:
			f.Close()
			return
		}
	}
}

func (f *File) sendClose() {
	f.closed <- true
}

func (f *File) Close() {
	if f.watcher != nil {
		f.watcher.Close()
		f.watcher = nil
	}

	if f.ch != nil {
		f.ch = nil
	}

	if f.fp != nil {
		f.fp.Close()
		f.fp = nil
	}
}

func (f *File) getCurrentFileSize() int64 {
	stat, err := f.fp.Stat()
	if err != nil {
	}
	return stat.Size()
}

func (f *File) openFile(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		panic("cant find file:" + filename)
	}
	f.fp = fp
	f.fp.Seek(0, os.SEEK_END)
}

func (f *File) watchFile(filename string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	f.watcher = watcher
	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}
}
