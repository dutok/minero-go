// Package direction is used with packet PlayerAction (0x0E).
package direction

// Value    0   1   2   3   4   5
// Offset   -Y  +Y  -Z  +Z  -X  +X
const (
	Up = iota
	Down
	Bottom
	Front
	Left
	Right
)
