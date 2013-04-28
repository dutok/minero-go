// Package tick a basic interface to work with tick events.
package tick

type Ticker interface {
	Tick(int64)
}
