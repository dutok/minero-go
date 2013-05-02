package world

import (
	"bytes"
	"compress/zlib"
	"github.com/toqueteos/minero/util/must"
	"io"
)

const (
	None = iota
	Gzip
	Zlib
)

// Region (.mca files) store a region's chunks, that's 32x32 chunks.
type Region struct {
	Pos  [1024]int32 // Chunk position in 4k increments from start.
	Mod  [1024]int32 // Last modification time of a chunk
	Data [1024]struct {
		Length int32
		Chunk  []byte
	}
}

func (re *Region) ReadFrom(r io.Reader) (n int64, err error) {
	var rw must.ReadWriter

	// Copy everything to a buffer. Max size: 4MB + 8KB
	var all bytes.Buffer
	rw.Must(io.Copy(&all, r))

	// Read chunk positions.
	for i := 0; i < len(re.Pos); i++ {
		// Read 4KB offset from file start. Only first 3 bytes needed.
		re.Pos[i] = rw.ReadInt32(&all) >> 8

		// Fourth byte is a 4KB section counter which is ignored because we
		// already know the length of chunk data.
		//
		// More info here:
		// http://www.minecraftwiki.net/wiki/Region_file_format#Structure
		//
		// " The remainder of the file consists of data for up to 1024 chunks,
		// interspersed with an arbitrary amount of unused space. "
		//
		// TLDR: Just another idiotic/bad designed spec.
	}

	// Read chunk timestamps.
	//
	// Last modification time of a chunk. Unit: unknown, seconds?
	//
	// NOTE: Does something use this? MCEdit maybe?
	for i := 0; i < len(re.Mod); i++ {
		re.Mod[i] = rw.ReadInt32(&all)
	}

	// Read chunk data.
	for i := 0; i < len(re.Data); i++ {
		re.Data[i].Length = rw.ReadInt32(&all)
		re.Data[i].Compression = byte(rw.ReadInt8(&all))

		var buf bytes.Buffer
		io.CopyN(&buf, &all, length-1)

		switch scheme {
		case Gzip:
			panic("Alpha chunk format not implemented.")
		case Zlib:
			zr := zlib.NewReader(&all)
			io.Copy(&buf, zr)
		}

		re.Data[i].Chunk = buf.Bytes()
	}

	return rw.Result()
}

func (r *Region) WriteTo(w io.Writer) (n int64, err error) { return }
func (r *Region) ChunkPos(x, z int32) int                  { return z<<5 + x }

type Chunk struct {
	X, Z             int32        // XZ position of the chunk
	LastUpdate       int64        // Tick when the chunk was last saved.
	TerrainPopulated bool         // false=NMC resets world.
	Biomes           [256]byte    // -1=NMC reset biome.
	SkyLight         [256]int32   // Lowest Y light is at full strength. ZX.
	Sections         [16]*Section // 16x16x16 blocks.
	Entities         []Entity     // List of NBT Compound.
	TileEntities     []TileEntity // List of NBT Compound.
	TileTicks        []TileTick   // List of NBT Compound.
}

// NBT lookups:
// Level
//     xPos
//     zPos
//     LastUpdate
//     TerrainPopulated
//     InhabitedTime
//     Biomes
//     HeightMap
//     Sections
//         Y
//         Blocks
//         Add
//         Data
//         BlockLight
//         SkyLight
//     Entities
//     TileEntities
//     TileTicks

func (c *Chunk) ReadFrom(r io.Reader) (n int64, err error) { return }
func (c *Chunk) WriteTo(w io.Writer) (n int64, err error)  { return }

// Section defines one 16x16x16 block area within a chunk. A chunk is made of 16 sections.
type Section struct {
	Y          byte       // Y section index. 0~15 bottom to top.
	Blocks     [4096]byte // 4096B. 8b/block. YZX.
	Add        [2048]byte // 2048B. 4b/block. YZX. Add << 8 | Blocks
	Data       [2048]byte // 2048B. 4b/block. YZX.
	BlockLight [2048]byte // 2048B. 4b/block. YZX.
	SkyLight   [2048]byte // 2048B. 4b/block. YZX.
}

type Entity struct {
	Id       string     // Entity ID. This tag does not exist for the Player entity.
	Pos      [3]float64 // 3 TAG_Doubles describing the current XYZ position of the entity.
	Motion   [3]float64 // 3 TAG_Doubles describing the current dX,dY,dZ velocity of the entity in meters per tick.
	Rotation [2]float32 // Two TAG_Floats representing rotation in degrees.
	// The entity's rotation clockwise around the Y axis (called yaw). Due west is 0. Can have large values because it accumulates all of the entity's lateral rotation throughout the game.
	// The entity's declination from the horizon (called pitch). Horizontal is 0. Positive values look downward. Does not exceed positive or negative 90 degrees.
	FallDistance   float32 // Distance the entity has fallen. Larger values cause more damage when the entity lands.
	Fire           int16   // Number of ticks until the fire is put out. Negative values reflect how long the entity can stand in fire before burning.
	Air            int16   // How much air the entity has, in ticks. Fills to a maximum of 200 in air, giving 10 seconds submerged before the entity starts to drown, and a total of up to 20 seconds before the entity dies. Decreases while underwater. If 0 while underwater, the entity loses 1 health per second.
	OnGround       bool    // 1 or 0 (true/false) - true if the entity is touching the ground.
	Dimension      int32   // Unknown usage; entities are only saved in the region files for the dimension they are in. -1 for The Nether, 0 for The Overworld, and 1 for The End.
	Invulnerable   bool    // 1 or 0 (true/false) - true if the entity should not take damage. This applies to living and nonliving entities alike: mobs will not take damage from any source (including potion effects) and objects such as vehicles and item frames cannot be destroyed unless their supports are removed. Note that these entities also cannot be moved by fishing rods, attacks, explosions, or projectiles.
	PortalCooldown int32   // The number of ticks before which the entity may be teleported back through a portal of any kind. Initially starts at 900 ticks (45 seconds) and counts down.
	UUIDLeast      int64   // The least significant bits of this entity's Universally Unique IDentifier. This is joined with UUIDMost to form this entity's unique ID, which is currently unused.
	UUIDMost       int64   // The most significant bits of this entity's Universally Unique IDentifier.
	Riding         *Entity // The data of the entity being ridden. See this format (recursive).
}

type TileEntity struct {
	Id      string // Tile entity ID
	X, Y, Z int32  // Position
}

type TileTick struct {
	I       int32 // The ID of the block; used to activate the correct block update procedure.
	T       int32 // Number of ticks until processing should occur. May be negative when processing is overdue.
	X, Y, Z int32 // Position
}
