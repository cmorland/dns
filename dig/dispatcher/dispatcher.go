package dispatcher

import (
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records"
	"github.com/jawr/dns/database/models/tlds"
	"github.com/jawr/dns/dig/worker"
	"github.com/jawr/dns/log"
)

var Workers chan chan worker.Request
var Work chan worker.Request

type Result chan []records.Record

func init() {
	Start(2)
}

func Start(nworkers int) {
	Workers = make(chan chan worker.Request, nworkers)
	Work = make(chan worker.Request, nworkers*10)

	for i := 0; i < nworkers; i++ {
		w := worker.New(i+1, Workers)
		w.Start()
	}

	go func() {
		for {
			select {
			case work := <-Work:
				w := <-Workers
				w <- work
			}
		}
	}()
}

func AddDomain(d domains.Domain) Result {
	res := make(Result)
	wr := worker.Request{
		Domain: d,
		Result: res,
	}
	Work <- wr
	return res
}

func AddQuery(q string) Result {
	res := make(Result)
	// could offload this in to a seperate anon function to avoid bottleneck
	s, t, err := tlds.DetectDomainAndTLD(q)
	if err != nil {
		log.Error("Dig dipatcher: Unable to detect TLD and domain: %s (%s)", err, q)
		return res
	}
	d, err := domains.GetByNameAndTLD(s, t).One()
	if err != nil {
		log.Error("Dig dispatcher: unable to get domain: %s (%s)", err, s)
		d = domains.New(s, t)
		err = d.Insert()
		if err != nil {
			log.Error("Dig dispatcher: unable to insert domain: %s (%s)", err, d.String())
			return res
		}
	}
	wr := worker.Request{
		Domain: d,
		Result: res,
	}
	Work <- wr
	return res
}

// TODO: add some sort of AddQueryTimeout function that waits for a response

// TODO: add a stop/quit
