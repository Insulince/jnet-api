package ws

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	Message struct {
		NetworkId string
		Data      string
	}

	Worker struct {
		Active           bool
		Stream           chan Message
		ConnectionClosed chan struct{}
		Conn             *websocket.Conn
		Id               int
	}
)

var (
	wi      int
	workers []*Worker
)

func NewWorker(conn *websocket.Conn) *Worker {
	var w Worker

	wi++

	w.Active = true
	w.Stream = make(chan Message)
	w.ConnectionClosed = make(chan struct{})
	w.Conn = conn
	w.Id = wi

	workers = append(workers, &w)

	return &w
}

func (w *Worker) SendMessage(data string) error {
	err := w.Conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		return errors.Wrap(err, "writing text message to connection")
	}
	return nil
}

func (w *Worker) Close() {
	w.Active = false
	close(w.Stream)
	close(w.ConnectionClosed)
	// TODO(justin): Ignored error.
	_ = w.Conn.Close()
}

func NotifyActiveWorkers(networkId string, data string) {
	// TODO(justin): This seems poorly designed.

	var activeWorkers []*Worker
	for _, w := range workers {
		if w.Active {
			activeWorkers = append(activeWorkers, w)

			var m Message
			m.NetworkId = networkId
			m.Data = data

			w.Stream <- m
		}
	}
	workers = activeWorkers
}
