package writer

import (
	"container/list"
	"io"
)

type Broadcaster struct {
	channelList *list.List
	readCh      chan []byte
	tasks       chan func()
	closed      chan bool
}

func NewBroadcaster(readCh chan []byte) *Broadcaster {
	return &Broadcaster{
		list.New(),
		readCh,
		make(chan func(), 100),
		make(chan bool)}
}

func (b *Broadcaster) Add(w io.Writer) {
	f := func() {
		b.channelList.PushBack(w)
	}
	b.tasks <- f
}

func (b *Broadcaster) Remove(w io.Writer) {
	f := func() {
		elem := b.find(w)
		if elem != nil {
			b.channelList.Remove(elem)
		}
	}
	b.tasks <- f
}

func (b *Broadcaster) Exists(w io.Writer) bool {
	elem := b.find(w)
	return elem != nil
}

func (b *Broadcaster) find(w io.Writer) *list.Element {
	for elem := b.channelList.Front(); elem != nil; elem = elem.Next() {
		if elem.Value == w {
			return elem
		}
	}

	return nil
}

func (b *Broadcaster) Write(data []byte) {
	f := func() {
		for elem := b.channelList.Front(); elem != nil; elem = elem.Next() {
			w := elem.Value.(io.Writer)
			_, err := w.Write(data)
			if err != nil {
				b.Remove(w)
			}
		}
	}
	b.tasks <- f
}

func (b *Broadcaster) Close() {
	b.closed <- true
	<-b.closed
}

func (b *Broadcaster) Run() {
	for {
		select {
		case line := <-b.readCh:
			b.Write(line)
		case f := <-b.tasks:
			f()
		case <-b.closed:
			b.closed <- true
			return
		}
	}
}
