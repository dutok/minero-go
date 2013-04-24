// Package id implementes an entity id generator.
//
// Thanks to #mcdevs from Freenode, especially TkTech for the idea.
package id

import (
	"math"
	"sync"
	"sync/atomic"
)

// Generator returns numbers, not necessarily sequential, within the range:
// [0, 1<<31-1].
type Generator struct {
	max int32
	queue
}

// Get gets a free id.
func (g *Generator) Get() int32 {
	if !g.queue.empty() {
		return g.queue.pop()
	}
	if g.max == math.MaxInt32 {
		// Reset counter to 0? Can this be even possible?
		atomic.StoreInt32(&g.max, 0)
		return atomic.LoadInt32(&g.max)
	}
	return atomic.AddInt32(&g.max, 1)
}

// Rel releases an id back to the generator.
func (g *Generator) Rel(id int32) {
	g.queue.push(id)
}

var DefaultGenerator = new(Generator)

func Get() int32   { return DefaultGenerator.Get() }
func Rel(id int32) { DefaultGenerator.Rel(id) }

// queue is a FIFO queue capped to 1<<16 elements.
type queue struct {
	sync.RWMutex
	store []int32
	len   uint16
}

func (s queue) empty() bool {
	s.RLock()
	defer s.RUnlock()
	return s.len == 0
}

func (s queue) full() bool {
	s.RLock()
	defer s.RUnlock()
	return s.len == math.MaxUint16
}

// push adds an id to queue's end, will panic if queue is full.
func (s queue) push(v int32) {
	if !s.full() {
		s.Lock()
		s.store = append(s.store, v)
		s.len++
		s.Unlock()
		return
	}
	panic("can't push on full queue")
}

// push pops deletes and returns an id from queue's start, will panic if queue
// is empty.
func (s queue) pop() (i int32) {
	if !s.empty() {
		s.Lock()
		// Fetch value on top
		i = s.store[0]
		s.store = s.store[1:]
		s.Unlock()
		return
	}
	panic("can't pop on empty queue")
}
