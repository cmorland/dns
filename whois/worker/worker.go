package worker

import (
	"github.com/jawr/dns/whois/parser"
)

type Worker struct {
	ID          int
	Work        chan Request
	WorkerQueue chan chan Request
	QuitChan    chan bool
	Parser      parser.Parser
}

func New(id int, workerQueue chan chan Request) Worker {
	return Worker{
		ID:          id,
		Work:        make(chan Request),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		Parser:      parser.New(),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			// register with dispatcher
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				work.Do(w)

			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
