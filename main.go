package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jksy/go-log-viewer/config"
	"github.com/jksy/go-log-viewer/reader"
	"github.com/jksy/go-log-viewer/writer"
)

var broadcaster *writer.Broadcaster

func main() {
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		panic("load config error:" + err.Error())
	}

	readCh := make(chan []byte, 200)
	prepareReader(conf, readCh)
	broadcaster = prepareWriter(readCh)

	setupWebsocketUrl("/log/", broadcaster)
	setupIndexPage()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func prepareReader(conf *config.Config, readCh chan []byte) {
	for _, v := range conf.Inputs {
		u, err := url.Parse(v)
		if err != nil {
			log.Printf("input url parse error:%s, %s", v, err.Error())
		}
		if u.Scheme != "file" {
			log.Printf("cant support [%s] scheme, %s", u.Scheme, v)
			continue
		}
		fileReader := reader.NewFileReaderWithOpen(u.Path, readCh)
		go fileReader.Run()
	}
}

func prepareWriter(readCh chan []byte) *writer.Broadcaster {
	bc := writer.NewBroadcaster(readCh)
	go bc.Run()
	return bc
}

func setupIndexPage() {
	http.Handle("/", http.FileServer(http.Dir("./public/")))
}

func setupWebsocketUrl(url string, bc *writer.Broadcaster) {
	if bc == nil {
		panic("broadcaster is nil")
	}

	broadcaster = bc
	http.Handle(url, websocket.Handler(newConnectionHandler))
}

func newConnectionHandler(ws *websocket.Conn) {
	broadcaster.Add(ws)
	defer ws.Close()

	for {
		select {
		case <-time.After(1 * time.Second):
			if broadcaster.Exists(ws) == false {
				log.Println("disconnected")
				return
			}
		}
	}
}
