package es

import (
	"time"

	"gopkg.in/olivere/elastic.v3"
)

const (
	maxNumRates = 64
)

// Request represents a bulk request and its generation time.
type Request struct {
	Bulk *elastic.BulkService
	Took uint64
}

// Equalizer represents an equalzier to apply backpressure to bulk requests.
type Equalizer struct {
	Send  chan Request
	Ready chan error
	size  int
	rates []uint64
}

// NewEqualizer creates and returns a new equalizer object of a specific size.
func NewEqualizer(size int) *Equalizer {
	return &Equalizer{
		Send:  make(chan Request),
		Ready: make(chan error),
		size:  size,
	}
}

func (e *Equalizer) getAvg() float64 {
	total := uint64(0)
	for _, ms := range e.rates {
		total += ms
	}
	return float64(total) / float64(len(e.rates))
}

func (e *Equalizer) measure(ms uint64) {
	e.rates = append(e.rates, ms)
	if len(e.rates) > maxNumRates {
		// if past max rates, pop oldest one off
		e.rates = e.rates[1:]
	}
}

func (e *Equalizer) throttle(took uint64) {
	// get difference between the time it took to generate the payload vs
	// the time it takes ES to ingest
	delta := e.getAvg() - float64(took)
	// wait the difference if it is positive
	if delta > 0 {
		time.Sleep(time.Millisecond * time.Duration(delta))
	}
}

func (e *Equalizer) forwardRequest(req Request) {
	e.throttle(req.Took)
	took, err := SendBulkRequest(req.Bulk)
	e.measure(took)
	e.Ready <- err
}

func (e *Equalizer) listenToReqs() {
	for req := range e.Send {
		go e.forwardRequest(req)
	}
}

// Listen starts the equalizer by sending a number of ready requests according
// to its configured size.
func (e *Equalizer) Listen() {
	go func() {
		// send as many ready messages as there are expected listeners
		for i := 0; i < e.size; i++ {
			e.Ready <- nil
		}
	}()
	go e.listenToReqs()
}

// Close disables the equalizer so that it no longer listens to any incoming bulk requests.
func (e *Equalizer) Close() {
	close(e.Ready)
	close(e.Send)
}
