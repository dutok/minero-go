package world

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/minero/minero/util/must"
)

const (
	None = iota
	Gzip
	Zlib
)

// Region (.mca files) store 32x32 chunks.
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

func (re *Region) WriteTo(w io.Writer) (n int64, err error) { return }
func (re *Region) ChunkPos(x, z int32) int                  { return z<<5 + x }

// Chunk
//
// Note: All fields are children of root's compound tag named "Level".
type Chunk struct {
	X                int32        `xPos`             // X position of the chunk
	Z                int32        `zPos`             // Z position of the chunk
	LastUpdate       int64        `LastUpdate`       // Tick when the chunk was last saved.
	TerrainPopulated bool         `TerrainPopulated` // false=NMC resets world.
	Biomes           [256]byte    `Biomes`           // -1=NMC reset biome.
	HeightMap        [256]int32   `HeightMap`        // Lowest Y light is at full strength. ZX.
	Sections         [16]*Section `Sections`         // 16x16x16 blocks.
	Entities         []Entity     `Entities`         // List of NBT Compound.
	TileEntities     []TileEntity `TileEntities`     // List of NBT Compound.
	TileTicks        []TileTick   `TileTicks`        // List of NBT Compound.
}

func (c *Chunk) ReadFrom(r io.Reader) (n int64, err error) { return }
func (c *Chunk) WriteTo(w io.Writer) (n int64, err error)  { return }

// Section holds the 1/16 part of a chunk (16x16x16).
type Section struct {
	Y          byte       // Y section index. 0~15 bottom to top.
	Blocks     [4096]byte // 8b/block. YZX.
	Add        [2048]byte // 4b/block. YZX. Add << 8 | Blocks
	Data       [2048]byte // 4b/block. YZX.
	BlockLight [2048]byte // 4b/block. YZX.
	SkyLight   [2048]byte // 4b/block. YZX.
}

func (s *Section) ReadFrom(r io.Reader) (n int64, err error) { return }
func (s *Section) WriteTo(w io.Writer) (n int64, err error)  { return }

type Entity struct {
	Id             string  // Entity ID. Doesn't exist for players.
	X, Y, Z        float64 // Pos.
	Dx, Dy, Dz     float64 // Velocity. Unit: meters per tick.
	Yaw, Pitch     float32 // Look. Unit: degrees.
	FallDistance   float32 // Distance the entity has fallen.
	Fire           int16   // Fire ticks left or inmune ticks iff Fire < 0.
	Air            int16   // Air ticks left. Max: 200 (10s). Decreases under water.
	OnGround       bool    // Captain Obvious!
	Dimension      int32   // Unknown usage.
	Invulnerable   bool    // Applies to living/nonliving entities.
	PortalCooldown int32   // Starts at 900 ticks (45s) and decrements.
	UUIDLeast      int64   // Unused.
	UUIDMost       int64   // Unused.
	Riding         *Entity // Entity being ridden. Recursive.
}

func (e *Entity) ReadFrom(r io.Reader) (n int64, err error) { return }
func (e *Entity) WriteTo(w io.Writer) (n int64, err error)  { return }

type TileEntity struct {
	Id      string // Tile entity Id.
	X, Y, Z int32  // Pos.
}

func (te *TileEntity) ReadFrom(r io.Reader) (n int64, err error) { return }
func (te *TileEntity) WriteTo(w io.Writer) (n int64, err error)  { return }

type TileTick struct {
	Id      int32 // Block Id.
	Ticks   int32 // Ticks until processing. Iff Ticks < 0: overdue.
	X, Y, Z int32 // Pos.
}

func (tt *TileTick) ReadFrom(r io.Reader) (n int64, err error) { return }
func (tt *TileTick) WriteTo(w io.Writer) (n int64, err error)  { return }
