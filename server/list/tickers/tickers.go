// Package tickers implements a goroutine-safe ticker list.
package tickers

import (
	"sync"

	"github.com/minero/minero/server/tick"
)

// Tickers is a simple goroutine-safe ticker list.
type Tickers struct {
	sync.RWMutex

	list map[int32]tick.Ticker
}

func New() Tickers {
	return Tickers{
		list: make(map[int32]tick.Ticker),
	}
}

// Len returns the number of active tickers.
func (l Tickers) Len() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.list)
}

// Copy returns a copy of the list.
func (l Tickers) Copy() map[int32]tick.Ticker {
	lc := make(map[int32]tick.Ticker)
	l.RLock()
	for k, t := range l.list {
		lc[k] = t
	}
	l.RUnlock()
	return lc
}

// GetTicker gets a ticker from the list by its Entity Id.
func (l Tickers) GetTicker(id int32) tick.Ticker {
	l.RLock()
	defer l.RUnlock()
	return l.list[id]
}

// AddTicker adds a ticker to the list.
func (l Tickers) AddTicker(id int32, t tick.Ticker) {
	l.Lock()
	l.list[id] = t
	l.Unlock()
}

// RemTicker removes a ticker from the list.
func (l Tickers) RemTicker(id int32) {
	l.Lock()
	delete(l.list, id)
	l.Unlock()
}

// TickAll calls each ticker's Tick method.
func (l Tickers) TickAll(tick int64) {
	l.RLock()
	for _, t := range l.list {
		t.Tick(tick)
	}
	l.RUnlock()
}
