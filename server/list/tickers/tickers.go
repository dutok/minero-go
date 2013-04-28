// Package tickers implements a goroutine-safe ticker list.
package tickers

import (
	"sync"

	"github.com/toqueteos/minero/server/tick"
)

// TickersList is a simple goroutine-safe ticker list.
type TickersList struct {
	sync.RWMutex

	list map[int32]tick.Ticker
}

func New() TickersList {
	return TickersList{
		list: make(map[int32]tick.Ticker),
	}
}

// Len returns the number of active tickers.
func (l TickersList) Len() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.list)
}

// Copy returns a copy of the list.
func (l TickersList) Copy() map[int32]tick.Ticker {
	lc := make(map[int32]tick.Ticker)
	l.RLock()
	for k, t := range l.list {
		lc[k] = t
	}
	l.RUnlock()
	return lc
}

// GetTicker gets a ticker from the list by its Entity Id.
func (l TickersList) GetTicker(id int32) tick.Ticker {
	l.RLock()
	defer l.RUnlock()
	return l.list[id]
}

// AddTicker adds a ticker to the list.
func (l TickersList) AddTicker(id int32, t tick.Ticker) {
	l.Lock()
	l.list[id] = t
	l.Unlock()
}

// RemTicker removes a ticker from the list.
func (l TickersList) RemTicker(id int32) {
	l.Lock()
	delete(l.list, id)
	l.Unlock()
}

// TickAll calls each ticker's Tick method.
func (l TickersList) TickAll(tick int64) {
	l.RLock()
	for _, t := range l.list {
		t.Tick(tick)
	}
	l.RUnlock()
}
