// Package dimension defines official Minecraft and custom dimensions.
package dimension

const (
	Nether = iota - 1
	Overworld
	End

	Custom = 1 << 7 // Unofficial
)
