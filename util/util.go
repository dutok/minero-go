// Package util defines utility functions.
package util

import (
	"github.com/toqueteos/minero/constants"
)

// Ticks translates s seconds to in-game ticks. Usually s * 20.
func Ticks(s int64) int64 {
	return s * constants.TicksPerSecond
}
