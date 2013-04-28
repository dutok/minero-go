// Package gamemode defines all available game modes and some utility functions.
package gamemode

const (
	Survival = iota
	Creative
	Adventure
)

// IsHardcore returns true if gamemode has the hardcore flag, false otherwise.
func IsHardcore(d int) bool {
	return d&0x8 == 0x8
}

// SetHardcore returns gamemode d with or without the hardcore flag.
func SetHardcore(d int, h bool) int {
	if h {
		return d | 0x8
	}
	return d & 0x3
}
