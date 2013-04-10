package constants

// Constants validation: http://play.golang.org/p/VwcPz3JoFh

type PlayerAction byte

const (
	StartDig PlayerAction = iota
	CancelDig
	FinishDig
	DropStack
	DropItem
	ShootArrow   = 5 // Same id
	FinishEating = 5 // Same id
)

// Source: http://wiki.vg/Block_Actions
type BlockAction int8

const (
	// Note Block
	// It shows the note particle being emitted from the block as well as playing the tone.
	// Byte 1: Instrument type.
	NoteBlockHarp BlockAction = iota
	NoteBlockDoubleBass
	NoteBlockSnareDrum
	NoteBlockClicksSticks
	NoteBlockBassDrum
	// Byte 2: Note pitch 0~24 (low-high). More information on Minecraft Wiki.

	// Piston
	// Byte 1: Piston state
	PistonPush BlockAction = iota
	PistonPull
	// Byte 2: Direction, see BlockDirection
	// PistonDirDown BlockAction = iota
	// PistonDirUp
	// PistonDirSouth
	// PistonDirWest
	// PistonDirNorth
	// PistonDirEast

	// Chest
	// Animates the chest's lid opening. Notchian server will send this every 3s
	// even if the state hasn't changed.
	// Byte 1: Not used. Always 1
	// Byte 2: State of the chest
	ChestClosed BlockAction = iota
	ChestOpen
)

type BlockDirection int8

const (
	// Value    0   1   2   3   4   5
	// Offset   -Y  +Y  -Z  +Z  -X  +X
	Up BlockDirection = iota
	Down
	Bottom
	Front
	Left
	Right
)

type Difficulty int8

const (
	Peaceful Difficulty = iota
	Easy
	Normal
	Hard
)

type Dimension int8

const (
	Nether Dimension = iota - 1
	Overworld
	End
)

type EntityStatus int8

const (
	_                    EntityStatus = iota // 0
	_                                        // 1
	EntityHurt           = 2
	EntityDead           = 3
	_                    // 4
	_                    // 5
	WolfTaming           = 6
	WolfTamed            = 7
	WolfShakingWater     = 8
	ServerAcceptedEating = 9
	SheepEating          = 10
)

type GameMode int8

const (
	Survival GameMode = iota
	Creative
	Adventure
)
