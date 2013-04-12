// Package action is used with packet BlockAction (0x36).
//
// Source: http://wiki.vg/Block_Actions
package action

// Note Block: It shows the note particle being emitted from the block as well
// as playing the tone.
// - Byte 1: Instrument type.
// - Byte 2: Note pitch 0~24 (low-high). More information on Minecraft Wiki.
const (
	NoteBlockHarp = iota
	NoteBlockDoubleBass
	NoteBlockSnareDrum
	NoteBlockClicksSticks
	NoteBlockBassDrum
)

// Piston
// - Byte 1: Piston state
// - Byte 2: Direction (constants/block/directionsee BlockDirection
const (
	PistonPush = iota
	PistonPull
)

const (
	PistonDirDown = iota
	PistonDirUp
	PistonDirSouth
	PistonDirWest
	PistonDirNorth
	PistonDirEast
)

// Chest: Animates the chest's lid opening. Notchian server will send this every
// 3s even if the state hasn't changed.
// - Byte 1: Not used. Always 1
// - Byte 2: State of the chest
const (
	ChestClosed = iota
	ChestOpen
)
