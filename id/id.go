// Package id
//
// Thanks to TkTech on #mdevs@freenode for the idea
package id

import (
	"sync/atomic"
)

var DefaultCounter Counter

// Counter and a capped stack.
// Counter is the max current entity ID, stack holds recently released IDs that can be reclaimed

type stack struct {
	store []int32
	len   int
}

func (s stack) empty() bool { return s.len == 0 }

func (s stack) push(v int32) {
	if s.len < len(s.store) {
		s.store[s.len] = v
	} else {
		s.store = append(s.store, v)
	}
	s.len++
}

func (s stack) pop() int32 {
	if s.len > 0 {
		// Fetch value on top
		v := s.store[s.len-1]
		// Delete old value
		s.store[s.len-1] = nil
		s.len--
		return v
	}
	return nil
}

func (s stack) top() int32 {
	if s.len > 0 {
		return s.store[s.len-1]
	}
	return nil
}

type Counter struct {
	c int32
}

func (c *Counter) New() int32 {
	atomic.LoadInt32(addr)
}

func New() int32 { return atomic.AddInt32(&c, 1) }
