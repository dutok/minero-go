package minecraft

// Iterate struct fields (type & value) with reflect:
// http://play.golang.org/p/BzYrOzevoJ

// Unit vector from pitch + yaw:
// http://play.golang.org/p/PSh5P13YMJ

// rw.MustRead([a-zA-Z]+)\(w,
// rw.MustWrite$1(w,

import (
	"encoding/binary"
	"fmt"
	"io"

	mct "github.com/toqueteos/minero/types/minecraft"
)

// Global io error handler
var rw PacketErrorHandler

type Packet interface {
	Id() byte
	io.ReaderFrom
	io.WriterTo
}

func CheckPacketId(expected, input byte) error {
	var en, in string
	var ok bool

	if en, ok = packetNames[expected]; !ok {
		en = "UnknownExpectedId"
	}
	if in, ok = packetNames[input]; !ok {
		in = "UnknownInputId"
	}

	if input != expected {
		return fmt.Errorf("%s (%02x) got %s (%02x)", en, expected, in, input)
	}

	return nil
}

// KeepAlive is a two-way packet.
//
// The server will frequently send out a keep-alive, each containing a random
// ID. The client must respond with the same packet.
//
// Total Size: 5 bytes
type KeepAlive struct {
	RandomId int32
}

func (p KeepAlive) Id() byte { return 0x00 }
func (p *KeepAlive) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.RandomId = rw.MustReadInt(r)

	return rw.Result()
}
func (p *KeepAlive) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.RandomId)

	return rw.Result()
}

// LoginRequest is a server to client packet.
//
// Total Size: 12 bytes + length of strings
type LoginRequest struct {
	Entity     int32
	LevelType  string // level-type in server.properties
	GameMode   int8   // 0: survival, 1: creative, 2: adventure. Bit 3 (0x8) is the hardcore flag
	Dimension  int8   // -1: nether, 0: overworld, 1: end
	Difficulty int8   // 0~3 for Peaceful, Easy, Normal, Hard
	// _WorldHeight int8 // NMS: always 0
	MaxPlayers int8 // Used by the client to draw the player list
}

func (p LoginRequest) Id() byte { return 0x01 }
func (p *LoginRequest) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.LevelType = rw.MustReadString(r)
	p.GameMode = rw.MustReadByte(r)
	p.Dimension = rw.MustReadByte(r)
	p.Difficulty = rw.MustReadByte(r)
	_ = rw.MustReadByte(r)
	p.MaxPlayers = rw.MustReadByte(r)

	return rw.Result()
}
func (p *LoginRequest) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.LevelType)
	rw.MustWriteByte(w, p.GameMode)
	rw.MustWriteByte(w, p.Dimension)
	rw.MustWriteByte(w, p.Difficulty)
	rw.MustWriteByte(w, 0) // see LoginRequest._WorldHeight
	rw.MustWriteByte(w, p.MaxPlayers)

	return rw.Result()
}

// Handshake is a client to server packet.
//
// Total Size: 10 bytes + length of strings
type Handshake struct {
	Version  int8
	Username string // Player attempting to connect
	Host     string
	Port     int32
}

func (p Handshake) Id() byte { return 0x02 }
func (p *Handshake) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Version = rw.MustReadByte(r)
	p.Username = rw.MustReadString(r)
	p.Host = rw.MustReadString(r)
	p.Port = rw.MustReadInt(r)

	return rw.Result()
}
func (p *Handshake) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Version)
	rw.MustWriteString(w, p.Username)
	rw.MustWriteString(w, p.Host)
	rw.MustWriteInt(w, p.Port)

	return rw.Result()
}

// ChatMessage is a two-way packet.
//
// The default server will check the message to see if it begins with a '/'. If
// it doesn't, the username of the sender is prepended and sent to all other
// clients (including the original sender). If it does, the server assumes it to
// be a command and attempts to process it. A message longer than 100 characters
// will cause the server to kick the client.
//
// Note: User input must be sanitized server-side
//
// Total Size: 3 bytes + length of strings
type ChatMessage struct {
	Message string
}

func (p ChatMessage) Id() byte { return 0x03 }
func (p *ChatMessage) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Message = rw.MustReadString(r)

	return rw.Result()
}
func (p *ChatMessage) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Message)

	return rw.Result()
}

// TimeUpdate is a server to client packet.
//
// Time is based on ticks, exactly 20 ticks per second. There are 24000 ticks in
// a day, making Minecraft days exactly 20 minutes long.
//
// NOTES: Time = Time + 20 % 24000
//
// Total Size: 17 Bytes
type TimeUpdate struct {
	WorldAge, Time int64
}

func (p TimeUpdate) Id() byte { return 0x04 }
func (p *TimeUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WorldAge = rw.MustReadLong(r)
	p.Time = rw.MustReadLong(r)

	return rw.Result()
}
func (p *TimeUpdate) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteLong(w, p.WorldAge)
	rw.MustWriteLong(w, p.Time)

	return rw.Result()
}

// EntityEquipment is a server to client packet.
// Total Size: 7 bytes + slot data
type EntityEquipment struct {
	Entity int32
	Slot   int16     // Equipment slot: 0=held, 1-4=armor slot
	Item   *mct.Slot // Item in slot format
}

func (p EntityEquipment) Id() byte { return 0x05 }
func (p *EntityEquipment) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Slot = rw.MustReadShort(r)
	// p.Item = rw.MustReadSlot(r)

	return rw.Result()
}
func (p *EntityEquipment) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteShort(w, p.Slot)
	// rw.MustWriteSlot(w, p.Item)

	return rw.Result()
}

// SpawnPosition is a server to client packet.
//
// Sent by the server after login to specify the coordinates of the spawn point.
// It can be sent at any time to update the point compasses point at.
//
// Total Size: 13 bytes
type SpawnPosition struct {
	X, Y, Z int32 // Spawn X, Y, Z in block coordinates
}

func (p SpawnPosition) Id() byte { return 0x06 }
func (p *SpawnPosition) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)

	return rw.Result()
}
func (p *SpawnPosition) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)

	return rw.Result()
}

// EntityInteract is a client to server packet.
//
// Sent when an entity attacks or right-clicks another entity.
//
// NMS: accept iff entity being attacked/used is visible without obstruction and
// within a 4-unit radius of the player's position.
//
// Total Size: 10 bytes
type EntityInteract struct {
	From, To    int32
	MouseButton bool // true=left click, false=right click
}

func (p EntityInteract) Id() byte { return 0x07 }
func (p *EntityInteract) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.From = rw.MustReadInt(r)
	p.To = rw.MustReadInt(r)
	p.MouseButton = rw.MustReadBool(r)

	return rw.Result()
}
func (p *EntityInteract) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.From)
	rw.MustWriteInt(w, p.To)
	rw.MustWriteBool(w, p.MouseButton)

	return rw.Result()
}

// HealthUpdate is a server to client packet.
//
// Updates a player's health. Starts at 5.0.
//
// Food won't decrease while saturation is over zero. Eating increases food and
// saturation.
//
// Total Size: 9 bytes
type HealthUpdate struct {
	Health     int16   // <=0: dead, 20: full HP
	Food       int16   // 0~20
	Saturation float32 // 0.0~5.0, integer increments
}

func (p HealthUpdate) Id() byte { return 0x08 }
func (p *HealthUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Health = rw.MustReadShort(r)
	p.Food = rw.MustReadShort(r)
	p.Saturation = rw.MustReadFloat(r)

	return rw.Result()
}
func (p *HealthUpdate) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, p.Health)
	rw.MustWriteShort(w, p.Food)
	rw.MustWriteFloat(w, p.Saturation)

	return rw.Result()
}

// Respawn is a server to client packet.
//
// NMC: 1 is always sent c->s
//
// Total Size: 11 bytes + length of string
type Respawn struct {
	Dimension   int32  // -1: The Nether, 0: The Overworld, 1: The End
	Difficulty  int8   // 0~3: Peaceful, Easy, Normal, Hard.
	GameMode    int8   // 0: survival, 1: creative, 2: adventure. Hardcore flag not included
	WorldHeight int16  // Defaults to 256
	LevelType   string // See 0x01 login
}

func (p Respawn) Id() byte { return 0x09 }
func (p *Respawn) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Dimension = rw.MustReadInt(r)
	p.Difficulty = rw.MustReadByte(r)
	p.GameMode = rw.MustReadByte(r)
	p.WorldHeight = rw.MustReadShort(r)
	p.LevelType = rw.MustReadString(r)

	return rw.Result()
}
func (p *Respawn) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Dimension)
	rw.MustWriteByte(w, p.Difficulty)
	rw.MustWriteByte(w, p.GameMode)
	rw.MustWriteShort(w, p.WorldHeight)
	rw.MustWriteString(w, p.LevelType)

	return rw.Result()
}

// Player is a client to server packet.
//
// Fall damage is applied when this state goes from false to true and player is
// 4+ blocks high from ground.
//
// Total Size: 2 bytes
type Player struct {
	OnGround bool // true if the client is on the ground, false otherwise
}

func (p Player) Id() byte { return 0x0A }
func (p *Player) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.OnGround = rw.MustReadBool(r)

	return rw.Result()
}
func (p *Player) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteBool(w, p.OnGround)

	return rw.Result()
}

// PlayerPos client to server packet.
//
//
// Updates the players XYZ position on the server.
//
// If Stance - Y is less than 0.1 or greater than 1.65, the stance is illegal
// and the client will be kicked with the message “Illegal Stance”.
//
// If the distance between the last known position of the player on the server
// and the new position set by this packet is greater than 100 units will result
// in the client being kicked for "You moved too quickly :( (Hacking?)".
//
// Also if the absolute number of X or Z is set greater than 3.2E7D the client
// will be kicked for "Illegal position".
//
// Total Size: 34 bytes
type PlayerPos struct {
	X, Y, Z  float64 // Absolute position
	Stance   float64 // Modifies players' bounding box while on stairs, crouching, etc...
	OnGround bool    // Same as 0x0A
}

func (p PlayerPos) Id() byte { return 0x0B }
func (p *PlayerPos) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadDouble(r)
	p.Y = rw.MustReadDouble(r)
	p.Stance = rw.MustReadDouble(r)
	p.Z = rw.MustReadDouble(r)
	p.OnGround = rw.MustReadBool(r)

	return rw.Result()
}
func (p *PlayerPos) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteDouble(w, p.X)
	rw.MustWriteDouble(w, p.Y)
	rw.MustWriteDouble(w, p.Stance)
	rw.MustWriteDouble(w, p.Z)
	rw.MustWriteBool(w, p.OnGround)

	return rw.Result()
}

// PlayerLook client to server packet.
//
//
// Updates the direction the player is looking in.
//
// Yaw is measured in degrees, and does not follow classical trigonometry rules.
// The unit circle of yaw on the xz-plane starts at (0, 1) and turns backwards
// towards (-1, 0), or in other words, it turns clockwise instead of
// counterclockwise. Additionally, yaw is not clamped to between 0 and 360
// degrees; any number is valid, including negative numbers and numbers greater
// than 360.
//
// Pitch is measured in degrees, where 0 is looking straight ahead, -90 is
// looking straight up, and 90 is looking straight down.
//
// Total Size: 10 bytes
type PlayerLook struct {
	Yaw, Pitch float32 // Absolute rotation on the XY Axis (degrees)
	OnGround   bool    // Same as 0x0A
}

func (p PlayerLook) Id() byte { return 0x0C }
func (p *PlayerLook) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Yaw = rw.MustReadFloat(r)
	p.Pitch = rw.MustReadFloat(r)
	p.OnGround = rw.MustReadBool(r)

	return rw.Result()
}
func (p *PlayerLook) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteFloat(w, p.Yaw)
	rw.MustWriteFloat(w, p.Pitch)
	rw.MustWriteBool(w, p.OnGround)

	return rw.Result()
}

// PlayerPosLook is a two-way packet.
//
// Note: When this packet is sent from the server, the Y and Stance fields are
// swapped.
//
// Total Size: 42 bytes
type PlayerPosLook struct {
	X, Y, Z    float64 // Absolute position
	Stance     float64 // Modifies players' bounding box while on stairs, crouching, etc...
	Yaw, Pitch float32 // Absolute rotation on the XY Axis (degrees)
	OnGround   bool    // Same as 0x0A
}

func (p PlayerPosLook) Id() byte { return 0x0D }
func (p *PlayerPosLook) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadDouble(r)
	p.Y = rw.MustReadDouble(r)
	p.Stance = rw.MustReadDouble(r)
	p.Z = rw.MustReadDouble(r)
	p.Yaw = rw.MustReadFloat(r)
	p.Pitch = rw.MustReadFloat(r)
	p.OnGround = rw.MustReadBool(r)

	return rw.Result()
}
func (p *PlayerPosLook) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteDouble(w, p.X)
	rw.MustWriteDouble(w, p.Y)
	rw.MustWriteDouble(w, p.Stance)
	rw.MustWriteDouble(w, p.Z)
	rw.MustWriteFloat(w, p.Yaw)
	rw.MustWriteFloat(w, p.Pitch)
	rw.MustWriteBool(w, p.OnGround)

	return rw.Result()
}

// PlayerAction client to server packet.
//
// NMS: accepts packet iff coordinates within a 6-unit radius from player's
// position.
//
// Actions:
// 0: start dig.
// 1: cancel dig.
// 2: finish dig.
// 3: drop full stack (all other values set to 0).
// 4: drop single item from stack (all other values set to 0).
// 5: shoot arrow / finish eating (face: 0xff, all other values set to 0).
//
// Total Size: 12 bytes
type PlayerAction struct {
	Action int8 // (see proto/minecraft/constants#PlayerAction)
	X      int32
	Y      int8
	Z      int32
	Face   int8 // (see proto/minecraft/constants#BlockDirection)
}

func (p PlayerAction) Id() byte { return 0x0E }
func (p *PlayerAction) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Action = rw.MustReadByte(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadInt(r)
	p.Face = rw.MustReadByte(r)

	return rw.Result()
}
func (p *PlayerAction) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Action)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Face)

	return rw.Result()
}

// PlayerBlockPlace client to server packet.
//
// Iff XYZ + Direction == -1 currently held item should have its state updated
// (eating food, shooting bows, using buckets, etc...)
//
// Note: NMC might send two packets when using buckets, first a normal and then
// a special case. First a normal packet is sent when you're looking at a block,
// it does nothing on NMS. Second packet performs the action (based on current
// pos/look and with a distance check, see next note).
//
// Note: buckets can only be used within a radius of 6 units.
//
// Total Size: 14 bytes + slot data
type PlayerBlockPlace struct {
	X             int32
	Y             byte
	Z             int32
	Direction     int8 // (see proto/minecraft/constants#BlockDirection)
	HeldItem      *mct.Slot
	ChX, ChY, ChZ int8 // The position of the crosshair on the block
}

func (p PlayerBlockPlace) Id() byte { return 0x0F }
func (p *PlayerBlockPlace) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = byte(rw.MustReadByte(r))
	p.Z = rw.MustReadInt(r)
	p.Direction = rw.MustReadByte(r)
	p.HeldItem = rw.MustReadSlot(r)
	p.ChX = rw.MustReadByte(r)
	p.ChY = rw.MustReadByte(r)
	p.ChZ = rw.MustReadByte(r)

	return rw.Result()
}
func (p *PlayerBlockPlace) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteByte(w, int8(p.Y))
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Direction)
	rw.MustWriteSlot(w, p.HeldItem)
	rw.MustWriteByte(w, p.ChX)
	rw.MustWriteByte(w, p.ChY)
	rw.MustWriteByte(w, p.ChZ)

	return rw.Result()
}

// ItemHeldChange is a two-way packet.
// Total Size: 3 bytes
type ItemHeldChange struct {
	SlotId int16 // The slot which the player has selected (0-8)
}

func (p ItemHeldChange) Id() byte { return 0x10 }
func (p *ItemHeldChange) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.SlotId = rw.MustReadShort(r)

	return rw.Result()
}
func (p *ItemHeldChange) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, p.SlotId)

	return rw.Result()
}

// BedUse is a server to client packet.
//
// Note: This Packet is sent to all nearby players including the one sent to bed.
//
// Total Size: 15 bytes
type BedUse struct {
	Entity int32
	X      int32
	Y      int8
	Z      int32
}

func (p BedUse) Id() byte { return 0x11 }
func (p *BedUse) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	_ = rw.MustReadByte(r) // Unknown use, only 0 observed
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadInt(r)

	return rw.Result()
}
func (p *BedUse) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, 0) // Unknown use, only 0 observed
	rw.MustWriteInt(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteInt(w, p.Z)

	return rw.Result()
}

// Animation is a two-way packet.
// Total Size: 6 bytes
type Animation struct {
	Entity    int32
	Animation int8 // (see proto/minecraft/constants#EntityAnimation)
}

func (p Animation) Id() byte { return 0x12 }
func (p *Animation) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Animation = rw.MustReadByte(r)

	return rw.Result()
}
func (p *Animation) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Animation)

	return rw.Result()
}

// EntityAction is a client to server packet.
// Total Size: 6 bytes
type EntityAction struct {
	Entity int32
	Action int8 // (see proto/minecraft/constants#EntityAction)
}

func (p EntityAction) Id() byte { return 0x13 }
func (p *EntityAction) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Action = rw.MustReadByte(r)

	return rw.Result()
}
func (p *EntityAction) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Action)

	return rw.Result()
}

// EntityNamedSpawn is a server to client packet.
//
// Note: sent when a player comes into visible range, not when it player joins.
//
// Note: Item <= 0 crashes clients.
//
// Total Size: 22 bytes + length of strings + metadata (at least 1)
type EntityNamedSpawn struct {
	Entity     int32
	Name       string // Max length: 16
	X, Y, Z    int32  // Absolute Integer
	Yaw, Pitch int8
	Item       int16 // Item currently holding. 0: no item
	Metadata   *mct.Metadata
}

func (p EntityNamedSpawn) Id() byte { return 0x14 }
func (p *EntityNamedSpawn) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Name = rw.MustReadString(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Yaw = rw.MustReadByte(r)
	p.Pitch = rw.MustReadByte(r)
	p.Item = rw.MustReadShort(r)
	p.Metadata = rw.MustReadMetadata(r)

	return rw.Result()
}
func (p *EntityNamedSpawn) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteString(w, p.Name)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteByte(w, p.Pitch)
	rw.MustWriteShort(w, p.Item)
	rw.MustWriteMetadata(w, p.Metadata)

	return rw.Result()
}

// ItemCollect is a server to client packet.
//
// Note: Server checks items to be picked up after each PlayerPos and
// PlayerPosLook packet sent.
//
// Total Size: 9 bytes
type ItemCollect struct {
	WhatId, WhoId int32
}

func (p ItemCollect) Id() byte { return 0x16 }
func (p *ItemCollect) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WhatId = rw.MustReadInt(r)
	p.WhoId = rw.MustReadInt(r)

	return rw.Result()
}
func (p *ItemCollect) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.WhatId)
	rw.MustWriteInt(w, p.WhoId)

	return rw.Result()
}

// SpawnObjectVehicle is a server to client packet.
// Total Size: 23 or 29 bytes
type SpawnObjectVehicle struct {
	Entity     int32
	Type       int8            // The type of object (see Entities#Objects)
	X, Y, Z    int32           // Absolute Integer Position of the object
	Pitch, Yaw int8            // In steps of 2p/256
	ObjectData *mct.ObjectData // (see Object Data)
}

func (p SpawnObjectVehicle) Id() byte { return 0x17 }
func (p *SpawnObjectVehicle) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Type = rw.MustReadByte(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Pitch = rw.MustReadByte(r)
	p.Yaw = rw.MustReadByte(r)
	p.ObjectData = rw.MustReadObjectData(r)

	return rw.Result()
}
func (p *SpawnObjectVehicle) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Type)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Pitch)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteObjectData(w, p.ObjectData)

	return rw.Result()
}

// SpawnMob is a server to client packet.
// Total Size: 27 bytes + Metadata (3+ bytes)
type SpawnMob struct {
	Entity                int32
	Type                  int8  // The type of object (see Entities#Objects)
	X, Y, Z               int32 // Absolute Integer Position of the object
	Pitch, HeadPitch, Yaw int8  // Yaw in steps of 2p/256
	VelX, VelY, VelZ      int16
	Metadata              *mct.Metadata // Varies by mob (see Entities)
}

func (p SpawnMob) Id() byte { return 0x18 }
func (p *SpawnMob) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Type = rw.MustReadByte(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Pitch = rw.MustReadByte(r)
	p.HeadPitch = rw.MustReadByte(r)
	p.Yaw = rw.MustReadByte(r)
	p.VelX = rw.MustReadShort(r)
	p.VelY = rw.MustReadShort(r)
	p.VelZ = rw.MustReadShort(r)
	p.Metadata = rw.MustReadMetadata(r)

	return rw.Result()
}

func (p *SpawnMob) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Type)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Pitch)
	rw.MustWriteByte(w, p.HeadPitch)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteShort(w, p.VelX)
	rw.MustWriteShort(w, p.VelY)
	rw.MustWriteShort(w, p.VelZ)
	rw.MustWriteMetadata(w, p.Metadata)

	return rw.Result()
}

// SpawnPainting is a server to client packet.
//
// Note: Title's max length is 13.
//
// Total Size: 23 bytes + length of string
type SpawnPainting struct {
	Entity    int32
	Title     string
	X, Y, Z   int32
	Direction int32 // Direction the painting faces (0: -z, 1: -x, 2: +z, 3: +x)
}

func (p SpawnPainting) Id() byte { return 0x19 }
func (p *SpawnPainting) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Title = rw.MustReadString(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Direction = rw.MustReadInt(r)

	return rw.Result()
}

func (p *SpawnPainting) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteString(w, p.Title)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteInt(w, p.Direction)

	return rw.Result()
}

// SpawnExperienceOrb is a server to client packet.
// Total Size: 19 bytes
type SpawnExperienceOrb struct {
	Entity  int32
	X, Y, Z int32 // Absolute
	Count   int16
}

func (p SpawnExperienceOrb) Id() byte { return 0x1A }
func (p *SpawnExperienceOrb) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Count = rw.MustReadShort(r)

	return rw.Result()
}

func (p *SpawnExperienceOrb) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteShort(w, p.Count)

	return rw.Result()
}

// EntityVelocity is a server to client packet.
//
// Velocity is believed to be in units of 1/8000 of a block per server tick
// (50ms).
//
// Example: -1343 would move (-1343 / 8000) = −0.167875 blocks per tick (or
// −3,3575 blocks per second).
//
// Total Size: 11 bytes
type EntityVelocity struct {
	Entity           int32
	VelX, VelY, VelZ int16 // (Protocol#Entity Velocity (0x1C))
}

func (p EntityVelocity) Id() byte { return 0x1C }
func (p *EntityVelocity) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.VelX = rw.MustReadShort(r)
	p.VelY = rw.MustReadShort(r)
	p.VelZ = rw.MustReadShort(r)

	return rw.Result()
}

func (p *EntityVelocity) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteShort(w, p.VelX)
	rw.MustWriteShort(w, p.VelY)
	rw.MustWriteShort(w, p.VelZ)

	return rw.Result()
}

// EntityDestroy is a server to client packet.
// Total Size: 2 + (entity count * 4) bytes
type EntityDestroy struct {
	Count    int8
	Entities []int32
}

func (p EntityDestroy) Id() byte { return 0x1D }
func (p *EntityDestroy) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Count = rw.MustReadByte(r)
	for i := 0; i < int(p.Count); i++ {
		p.Entities = append(p.Entities, rw.MustReadInt(r))
	}

	return rw.Result()
}

func (p *EntityDestroy) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, int8(len(p.Entities)))
	for index, _ := range p.Entities {
		rw.MustWriteInt(w, p.Entities[index])
	}

	return rw.Result()
}

// Entity is a server to client packet.
//
// Sent every game tick.
//
// Entity did not move/look since the last PlayerPos/Look packet.
//
// Total Size: 5 bytes
type Entity struct {
	Entity int32
}

func (p Entity) Id() byte { return 0x1E }
func (p *Entity) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)

	return rw.Result()
}

func (p *Entity) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)

	return rw.Result()
}

// EntityRelMove is a server to client packet.
//
// Sent when an entity moves less than 4 blocks, otherwise use 0x22.
//
// Total Size: 8 bytes
type EntityRelMove struct {
	Entity  int32
	X, Y, Z int8 // Axis Relative movement as an Absolute Integer
}

func (p EntityRelMove) Id() byte { return 0x1F }
func (p *EntityRelMove) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.X = rw.MustReadByte(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityRelMove) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteByte(w, p.Z)

	return rw.Result()
}

// EntityLook is a server to client packet.
//
// Sent when an entity rotates.
//
// Total Size: 7 bytes
type EntityLook struct {
	Entity     int32
	Yaw, Pitch int8 // The X & Y Axis rotation as a fraction of 360
}

func (p EntityLook) Id() byte { return 0x20 }
func (p *EntityLook) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Yaw = rw.MustReadByte(r)
	p.Pitch = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityLook) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteByte(w, p.Pitch)

	return rw.Result()
}

// EntityLookRelMove is a server to client packet.
//
// Mix of 0x1F + 0x20.
//
// Total Size: 10 bytes
type EntityLookRelMove struct {
	Entity     int32
	X, Y, Z    int8 // Axis Relative movement as an Absolute Integer
	Yaw, Pitch int8 // The X & Y Axis rotation as a fraction of 360
}

func (p EntityLookRelMove) Id() byte { return 0x21 }
func (p *EntityLookRelMove) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.X = rw.MustReadByte(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadByte(r)
	p.Yaw = rw.MustReadByte(r)
	p.Pitch = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityLookRelMove) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteByte(w, p.Z)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteByte(w, p.Pitch)

	return rw.Result()
}

// EntityTeleport is a server to client packet.
//
// Complementary of 0x1F. Sent when an entity moves more than 4 blocks.
//
// Total Size: 19 bytes
type EntityTeleport struct {
	Entity     int32
	X, Y, Z    int32 // Position as an Absolute Integer
	Yaw, Pitch int8  // The X & Y Axis rotation as a fraction of 360
}

func (p EntityTeleport) Id() byte { return 0x22 }
func (p *EntityTeleport) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Yaw = rw.MustReadByte(r)
	p.Pitch = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityTeleport) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Yaw)
	rw.MustWriteByte(w, p.Pitch)

	return rw.Result()
}

// EntityHeadLook is a server to client packet.
// Total Size: 6 bytes
type EntityHeadLook struct {
	Entity  int32
	HeadYaw int8 // Head yaw in steps of 2p/256
}

func (p EntityHeadLook) Id() byte { return 0x23 }
func (p *EntityHeadLook) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.HeadYaw = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityHeadLook) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.HeadYaw)

	return rw.Result()
}

// EntityStatus is a server to client packet.
// Total Size: 6 bytes
type EntityStatus struct {
	Entity int32
	Status int8 // (see proto/minecraft/constants#EntityStatus)
}

func (p EntityStatus) Id() byte { return 0x26 }
func (p *EntityStatus) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Status = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityStatus) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Status)

	return rw.Result()
}

// EntityAttach is a server to client packet.
//
// Sent when an entity is attached to an entity (Minecart).
//
// Total Size: 9 bytes
type EntityAttach struct {
	Entity    int32
	VehicleId int32 // The vehicle entity Id attached to (-1 for unattaching)
}

func (p EntityAttach) Id() byte { return 0x27 }
func (p *EntityAttach) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.VehicleId = rw.MustReadInt(r)

	return rw.Result()
}

func (p *EntityAttach) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteInt(w, p.VehicleId)

	return rw.Result()
}

// EntityMetadata is a server to client packet.
// Total Size: 5 bytes + Metadata
type EntityMetadata struct {
	Entity int32
	Meta   *mct.Metadata // (see Entities)
}

func (p EntityMetadata) Id() byte { return 0x28 }
func (p *EntityMetadata) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Meta = rw.MustReadMetadata(r)

	return rw.Result()
}

func (p *EntityMetadata) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteMetadata(w, p.Meta)

	return rw.Result()
}

// EntityEffect is a server to client packet.
// Total Size: 9 bytes
type EntityEffect struct {
	Entity    int32
	Effect    int8 // (see PotionEffect)
	Amplifier int8
	Duration  int16
}

func (p EntityEffect) Id() byte { return 0x29 }
func (p *EntityEffect) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Effect = rw.MustReadByte(r)
	p.Amplifier = rw.MustReadByte(r)
	p.Duration = rw.MustReadShort(r)

	return rw.Result()
}

func (p *EntityEffect) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Effect)
	rw.MustWriteByte(w, p.Amplifier)
	rw.MustWriteShort(w, p.Duration)

	return rw.Result()
}

// EntityEffectRemove is a server to client packet.
// Total Size: 6 bytes
type EntityEffectRemove struct {
	Entity int32
	Effect int8 // (see PotionEffect)
}

func (p EntityEffectRemove) Id() byte { return 0x2A }
func (p *EntityEffectRemove) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Effect = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EntityEffectRemove) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Effect)

	return rw.Result()
}

// SetExperience is a server to client packet.
// Total Size: 9 bytes
type SetExperience struct {
	Xp      float32
	Level   int16
	TotalXp int16
}

func (p SetExperience) Id() byte { return 0x2B }
func (p *SetExperience) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Xp = rw.MustReadFloat(r)
	p.Level = rw.MustReadShort(r)
	p.TotalXp = rw.MustReadShort(r)

	return rw.Result()
}

func (p *SetExperience) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteFloat(w, p.Xp)
	rw.MustWriteShort(w, p.Level)
	rw.MustWriteShort(w, p.TotalXp)

	return rw.Result()
}

// ChunkData is a server to client packet.
// Total Size: 18 bytes + len(ChunkData)
type ChunkData struct {
	X, Z               int32  // Chunk XZ Coordinate
	GroundUpContinuous bool   // true = all sections in this vertical column, where the primary bitmap specifies exactly which sections are included, and which are air.
	Primary            uint16 // Bitmask. 1 for every 16x16x16 section
	Add                uint16 // Bitmask. 1 for every 16x16x16 section ("add" on payload)
	// BUG(toqueteos): type should be proto/anvil/ChunkData instead of []byte
	ChunkData []byte // ZLib Deflate compressed chunk data
}

func (p ChunkData) Id() byte { return 0x33 }
func (p *ChunkData) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.GroundUpContinuous = rw.MustReadBool(r)
	p.Primary = uint16(rw.MustReadShort(r))
	p.Add = uint16(rw.MustReadShort(r))
	length := int(rw.MustReadInt(r))
	p.ChunkData = rw.MustReadByteArray(r, length)

	// // ChunkData is sent compressed with zlib deflate
	// var buf bytes.Buffer
	// zr := zlib.NewReader(bytes.NewBuffer(p.ChunkData))
	// rw.Must(io.Copy(&buf, zr))
	// p.ChunkData = buf.Bytes()

	return rw.Result()
}

func (p *ChunkData) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteBool(w, p.GroundUpContinuous)
	rw.MustWriteShort(w, int16(p.Primary))
	rw.MustWriteShort(w, int16(p.Add))
	rw.MustWriteInt(w, int32(len(p.ChunkData)))
	rw.MustWriteByteArray(w, p.ChunkData)

	// // ChunkData is sent compressed with zlib deflate
	// var buf bytes.Buffer
	// zr := zlib.NewWriter(bytes.NewBuffer(p.ChunkData))
	// rw.Must(io.Copy(&buf, zr))
	// length := int32(buf.Len())

	// rw.MustWriteInt(w, length)
	// rw.MustWriteByteArray(r, buf.Bytes())

	return rw.Result()
}

// BlockChangeMulti is a server to client packet.
// Total Size: 15 bytes + arrays
type BlockChangeMulti struct {
	X, Z   int32   // Chunk XZ Coordinate
	Count  int16   // len(Blocks)
	Length int32   // Length of payload
	Blocks []int32 // Coordinates, Type, and Metadata of blocks to change
}

func (p BlockChangeMulti) Id() byte { return 0x34 }
func (p *BlockChangeMulti) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Count = rw.MustReadShort(r)
	p.Length = rw.MustReadInt(r)
	p.Blocks = make([]int32, p.Count)
	for i := 0; i < int(p.Count); i++ {
		p.Blocks[i] = rw.MustReadInt(r)
	}

	return rw.Result()
}

func (p *BlockChangeMulti) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteShort(w, p.Count)
	rw.MustWriteInt(w, p.Length)
	for i := 0; i < len(p.Blocks); i++ {
		rw.MustWriteInt(w, p.Blocks[i])
	}

	return rw.Result()
}

// BlockChange is a server to client packet.
// Total Size: 13 bytes
type BlockChange struct {
	X         int32
	Y         int8
	Z         int32
	BlockType int16 // New block type for block
	BlockMeta int8  // New Metadata for block
}

func (p BlockChange) Id() byte { return 0x35 }
func (p *BlockChange) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadInt(r)
	p.BlockType = rw.MustReadShort(r)
	p.BlockMeta = rw.MustReadByte(r)

	return rw.Result()
}

func (p *BlockChange) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteShort(w, p.BlockType)
	rw.MustWriteByte(w, p.BlockMeta)

	return rw.Result()
}

// BlockAction is a server to client packet.
//
// It is used for:
// - Chests opening and closing
// - Pistons pushing and pulling
// - Note blocks playing
//
// Total Size: 15 bytes
type BlockAction struct {
	X            int32
	Y            int16
	Z            int32
	Byte1, Byte2 int8  // Varies depending on block (see Block_Actions)
	BlockId      int16 // The block id this action is set for
}

func (p BlockAction) Id() byte { return 0x36 }
func (p *BlockAction) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadShort(r)
	p.Z = rw.MustReadInt(r)
	p.Byte1 = rw.MustReadByte(r)
	p.Byte2 = rw.MustReadByte(r)
	p.BlockId = rw.MustReadShort(r)

	return rw.Result()
}

func (p *BlockAction) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteShort(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Byte1)
	rw.MustWriteByte(w, p.Byte2)
	rw.MustWriteShort(w, p.BlockId)

	return rw.Result()
}

// BlockBreakAnimation is a server to client packet.
// Total Size: 18 bytes
type BlockBreakAnimation struct {
	Entity  int32
	X, Y, Z int32 // Block position
	Damage  int8  // How far destroyed this block is
}

func (p BlockBreakAnimation) Id() byte { return 0x37 }
func (p *BlockBreakAnimation) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Damage = rw.MustReadByte(r)

	return rw.Result()
}

func (p *BlockBreakAnimation) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Damage)

	return rw.Result()
}

// MapChunkBulk is a server to client packet.
// Total Size: 8 + (DataLength) + 12 * (Count) bytes
type MapChunkBulk struct {
	Count        int16  // len(Meta)
	Length       int32  // len(ChunkData)
	SkylightSent bool   // Chunk data contains a light nibble array? true for overworld, false otherwise
	ChunkData    []byte // Compressed chunk data
	ChunkMeta    []ChunkMeta
}

type ChunkMeta struct {
	X, Z    int32  // The XZ coordinate of the specific chunk
	Primary uint16 // Bitmap. Specifies which sections are not empty
	Add     uint16 // Bitmap. Specifies which sections need add information because of very high block ids
}

func (p MapChunkBulk) Id() byte { return 0x38 }
func (p *MapChunkBulk) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Count = rw.MustReadShort(r)
	p.Length = rw.MustReadInt(r)
	p.SkylightSent = rw.MustReadBool(r)
	p.ChunkData = rw.MustReadByteArray(r, int(p.Length))
	p.ChunkMeta = make([]ChunkMeta, p.Length)
	for i := 0; i < int(p.Count); i++ {
		p.ChunkMeta[i].X = rw.MustReadInt(r)
		p.ChunkMeta[i].Z = rw.MustReadInt(r)
		p.ChunkMeta[i].Primary = uint16(rw.MustReadShort(r))
		p.ChunkMeta[i].Add = uint16(rw.MustReadShort(r))
	}

	return rw.Result()
}

func (p *MapChunkBulk) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, p.Count)
	rw.MustWriteInt(w, p.Length)
	rw.MustWriteBool(w, p.SkylightSent)
	rw.MustWriteByteArray(w, p.ChunkData)
	for i := 0; i < int(p.Count); i++ {
		rw.MustWriteInt(w, p.ChunkMeta[i].X)
		rw.MustWriteInt(w, p.ChunkMeta[i].Z)
		rw.MustWriteShort(w, int16(p.ChunkMeta[i].Primary))
		rw.MustWriteShort(w, int16(p.ChunkMeta[i].Add))
	}

	return rw.Result()
}

// Explosion is a server to client packet.
// Total Size: 45 bytes + 3*(Record count) bytes
type Explosion struct {
	X, Y, Z                   float64
	Radius                    float32 // NMC: Unused
	Count                     int32
	Blocks                    [][3]int8 // len(BlockPos) == Blocks * 3
	PlayerX, PlayerY, PlayerZ float32   // XYZ velocity of the player being pushed
}

func (p Explosion) Id() byte { return 0x3C }
func (p *Explosion) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadDouble(r)
	p.Y = rw.MustReadDouble(r)
	p.Z = rw.MustReadDouble(r)
	p.Radius = rw.MustReadFloat(r)
	p.Count = rw.MustReadInt(r)
	p.Blocks = make([][3]int8, p.Count)
	for i := 0; i < int(p.Count); i++ {
		p.Blocks[i] = [3]int8{
			rw.MustReadByte(r),
			rw.MustReadByte(r),
			rw.MustReadByte(r),
		}
	}
	p.PlayerX = rw.MustReadFloat(r)
	p.PlayerY = rw.MustReadFloat(r)
	p.PlayerZ = rw.MustReadFloat(r)

	return rw.Result()
}

func (p *Explosion) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteDouble(w, p.X)
	rw.MustWriteDouble(w, p.Y)
	rw.MustWriteDouble(w, p.Z)
	rw.MustWriteFloat(w, p.Radius)
	rw.MustWriteInt(w, p.Count)
	for i := 0; i < len(p.Blocks); i++ {
		rw.MustWriteByte(w, p.Blocks[i][0])
		rw.MustWriteByte(w, p.Blocks[i][1])
		rw.MustWriteByte(w, p.Blocks[i][2])
	}
	rw.MustWriteFloat(w, p.PlayerX)
	rw.MustWriteFloat(w, p.PlayerY)
	rw.MustWriteFloat(w, p.PlayerZ)

	return rw.Result()
}

// SoundEffect is a server to client packet.
//
// NMC: volume of sound effects is adjusted based on distance.
//
// Total Size: 19 bytes
type SoundEffect struct {
	EffectId         int32
	X                int32
	Y                int8
	Z                int32
	Data             int32 // Extra data for certain effects, see below.
	DisableRelVolume bool  // Ignored by client for all but mob.wither.spawn (1013)
}

func (p SoundEffect) Id() byte { return 0x3D }
func (p *SoundEffect) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.EffectId = rw.MustReadInt(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadByte(r)
	p.Z = rw.MustReadInt(r)
	p.Data = rw.MustReadInt(r)
	p.DisableRelVolume = rw.MustReadBool(r)

	return rw.Result()
}

func (p *SoundEffect) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.EffectId)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteByte(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteInt(w, p.Data)
	rw.MustWriteBool(w, p.DisableRelVolume)

	return rw.Result()
}

// SoundEffectNamed is a server to client packet.
// Total Size: 20 bytes + length of string
type SoundEffectNamed struct {
	Name    string  // 250
	X, Y, Z int32   // XYZ, multiplied by 8
	Volume  float32 // 1 is 100%, can be more
	Pitch   int8    // 63 is 100%, can be more
}

func (p SoundEffectNamed) Id() byte { return 0x3E }
func (p *SoundEffectNamed) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)
	p.Volume = rw.MustReadFloat(r)
	p.Pitch = rw.MustReadByte(r)

	return rw.Result()
}

func (p *SoundEffectNamed) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteFloat(w, p.Volume)
	rw.MustWriteByte(w, p.Pitch)

	return rw.Result()
}

// Particle is a server to client packet.
// Total Size: 34 bytes + length of string
type Particle struct {
	Name                      string
	X, Y, Z                   float32
	OffsetX, OffsetY, OffsetZ float32 // Added to XYZ after multiplication by random.nextGaussian()
	Speed                     float32
	Number                    int32
}

func (p Particle) Id() byte { return 0x3F }
func (p *Particle) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	p.X = rw.MustReadFloat(r)
	p.Y = rw.MustReadFloat(r)
	p.Z = rw.MustReadFloat(r)
	p.OffsetX = rw.MustReadFloat(r)
	p.OffsetY = rw.MustReadFloat(r)
	p.OffsetZ = rw.MustReadFloat(r)
	p.Speed = rw.MustReadFloat(r)
	p.Number = rw.MustReadInt(r)

	return rw.Result()
}

func (p *Particle) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteFloat(w, p.X)
	rw.MustWriteFloat(w, p.Y)
	rw.MustWriteFloat(w, p.Z)
	rw.MustWriteFloat(w, p.OffsetX)
	rw.MustWriteFloat(w, p.OffsetY)
	rw.MustWriteFloat(w, p.OffsetZ)
	rw.MustWriteFloat(w, p.Speed)
	rw.MustWriteInt(w, p.Number)

	return rw.Result()
}

// GameStateChange is a server to client packet.
// Total Size: 3 bytes
type GameStateChange struct {
	Reason   int8
	GameMode int8 // Reason == 3 then 0=survival, 1=creative
}

func (p GameStateChange) Id() byte { return 0x46 }
func (p *GameStateChange) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Reason = rw.MustReadByte(r)
	p.GameMode = rw.MustReadByte(r)

	return rw.Result()
}

func (p *GameStateChange) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Reason)
	rw.MustWriteByte(w, p.GameMode)

	return rw.Result()
}

// EntityGlobalSpawn is a server to client packet.
// Total Size: 18 bytes
type EntityGlobalSpawn struct {
	Entity  int32
	Type    int8  // Global entity type, currently always 1 for thunderbolt.
	X, Y, Z int32 // Thunderbolt XYZ as Absolute Integer
}

func (p EntityGlobalSpawn) Id() byte { return 0x47 }
func (p *EntityGlobalSpawn) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Entity = rw.MustReadInt(r)
	p.Type = rw.MustReadByte(r)
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadInt(r)
	p.Z = rw.MustReadInt(r)

	return rw.Result()
}

func (p *EntityGlobalSpawn) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.Entity)
	rw.MustWriteByte(w, p.Type)
	rw.MustWriteInt(w, p.X)
	rw.MustWriteInt(w, p.Y)
	rw.MustWriteInt(w, p.Z)

	return rw.Result()
}

// WindowOpen is a server to client packet.
// Total Size: 7 bytes + length of string
type WindowOpen struct {
	WindowId         int8   // Unique Window Id number. NMS: counter, starting at 1.
	InventoryType    int8   // The window type to use for display. Check below
	WindowTitle      string // The title of the window.
	Slots            int8   // Slots in window excluding player inventory slots
	UseProvidedTitle bool   // true: client uses what the server provides, false: client will look up a string
}

func (p WindowOpen) Id() byte { return 0x64 }
func (p *WindowOpen) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.InventoryType = rw.MustReadByte(r)
	p.WindowTitle = rw.MustReadString(r)
	p.Slots = rw.MustReadByte(r)
	p.UseProvidedTitle = rw.MustReadBool(r)

	return rw.Result()
}

func (p *WindowOpen) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteByte(w, p.InventoryType)
	rw.MustWriteString(w, p.WindowTitle)
	rw.MustWriteByte(w, p.Slots)
	rw.MustWriteBool(w, p.UseProvidedTitle)

	return rw.Result()
}

// WindowClose is a two-way packet.
// Total Size: 2 bytes
type WindowClose struct {
	WindowId int8 // 0 = inventory
}

func (p WindowClose) Id() byte { return 0x65 }
func (p *WindowClose) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)

	return rw.Result()
}

func (p *WindowClose) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)

	return rw.Result()
}

// WindowClick is a client to server packet.
// Total Size: 8 bytes + slot data
type WindowClick struct {
	WindowId int8      // 0: inventory
	Slot     int16     //
	Button   int8      // 0: left, 1: right, 3: middle ("Mode" set to 3)
	Action   int16     // Unique number for the action, used for transaction handling
	Mode     int8      // 0: regular, 1: shift + click, 5: "painting" mode, 6: double-click
	Item     *mct.Slot // Clicked item
}

func (p WindowClick) Id() byte { return 0x66 }
func (p *WindowClick) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Slot = rw.MustReadShort(r)
	p.Button = rw.MustReadByte(r)
	p.Action = rw.MustReadShort(r)
	p.Mode = rw.MustReadByte(r)
	p.Item = rw.MustReadSlot(r)

	return rw.Result()
}

func (p *WindowClick) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteShort(w, p.Slot)
	rw.MustWriteByte(w, p.Button)
	rw.MustWriteShort(w, p.Action)
	rw.MustWriteByte(w, p.Mode)
	rw.MustWriteSlot(w, p.Item)

	return rw.Result()
}

// WindowSlotSet is a server to client packet.
// Total Size: 4 bytes + slot data
type WindowSlotSet struct {
	WindowId int8 // 0 = inventory
	Slot     int16
	Data     *mct.Slot
}

func (p WindowSlotSet) Id() byte { return 0x67 }
func (p *WindowSlotSet) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Slot = rw.MustReadShort(r)
	p.Data = rw.MustReadSlot(r)

	return rw.Result()
}

func (p *WindowSlotSet) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteShort(w, p.Slot)
	rw.MustWriteSlot(w, p.Data)

	return rw.Result()
}

// WindowSetItems is a server to client packet.
// Total Size: 4 bytes + size of slot data array
type WindowSetItems struct {
	WindowId int8 // 0 = inventory
	Count    int16
	SlotData []*mct.Slot
}

func (p WindowSetItems) Id() byte { return 0x68 }
func (p *WindowSetItems) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Count = rw.MustReadShort(r)
	p.SlotData = make([]*mct.Slot, p.Count)
	for i := 0; i < int(p.Count); i++ {
		p.SlotData[i] = rw.MustReadSlot(r)
	}

	return rw.Result()
}

func (p *WindowSetItems) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteShort(w, p.Count)
	for i := 0; i < int(p.Count); i++ {
		rw.MustWriteSlot(w, p.SlotData[i])
	}

	return rw.Result()
}

// WindowUpdateProperty is a server to client packet.
// Total Size: 6 bytes
type WindowUpdateProperty struct {
	WindowId int8 // 0 = inventory
	Property int16
	Value    int16
}

func (p WindowUpdateProperty) Id() byte { return 0x69 }
func (p *WindowUpdateProperty) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Property = rw.MustReadShort(r)
	p.Value = rw.MustReadShort(r)

	return rw.Result()
}

func (p *WindowUpdateProperty) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteShort(w, p.Property)
	rw.MustWriteShort(w, p.Value)

	return rw.Result()
}

// ConfirmTransaction is a two-way packet.
// Total Size: 5 bytes
type ConfirmTransaction struct {
	WindowId int8
	Action   int16 // Unique number
	Accepted bool
}

func (p ConfirmTransaction) Id() byte { return 0x6A }
func (p *ConfirmTransaction) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Action = rw.MustReadShort(r)
	p.Accepted = rw.MustReadBool(r)

	return rw.Result()
}

func (p *ConfirmTransaction) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteShort(w, p.Action)
	rw.MustWriteBool(w, p.Accepted)

	return rw.Result()
}

// CreativeInventoryAction is a two-way packet.
// Total Size: 3 bytes + slot data
type CreativeInventoryAction struct {
	Slot int16
	Item *mct.Slot
}

func (p CreativeInventoryAction) Id() byte { return 0x6B }
func (p *CreativeInventoryAction) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Slot = rw.MustReadShort(r)
	p.Item = rw.MustReadSlot(r)

	return rw.Result()
}

func (p *CreativeInventoryAction) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, p.Slot)
	rw.MustWriteSlot(w, p.Item)

	return rw.Result()
}

// EnchantItem is a client to server packet.
//
// While the user is in the standard inventory (i.e., not a crafting bench) on a
// creative-mode server then the server will send this packet:
// - If an item is dropped into the quick bar.
// - If an item is picked up from the quick bar (item id is -1)..
//
// Total Size: 3 bytes
type EnchantItem struct {
	WindowId int8
	Position int8 // 0~2 from top to bottom
}

func (p EnchantItem) Id() byte { return 0x6C }
func (p *EnchantItem) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.WindowId = rw.MustReadByte(r)
	p.Position = rw.MustReadByte(r)

	return rw.Result()
}

func (p *EnchantItem) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.WindowId)
	rw.MustWriteByte(w, p.Position)

	return rw.Result()
}

// SignUpdate is a two-way packet.
//
// S->C: sent whenever a sign is discovered or created.
// C->S: sent when the "Done" button is pushed after placing a sign.
//
// Note: Not sent when a sign is destroyed or unloaded.
//
// Total Size: 11 bytes + 4 strings
type SignUpdate struct {
	X     int32
	Y     int16
	Z     int32
	Lines [4]string
}

func (p SignUpdate) Id() byte { return 0x82 }
func (p *SignUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadShort(r)
	p.Z = rw.MustReadInt(r)
	p.Lines[0] = rw.MustReadString(r)
	p.Lines[1] = rw.MustReadString(r)
	p.Lines[2] = rw.MustReadString(r)
	p.Lines[3] = rw.MustReadString(r)

	return rw.Result()
}

func (p *SignUpdate) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteShort(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteString(w, p.Lines[0])
	rw.MustWriteString(w, p.Lines[1])
	rw.MustWriteString(w, p.Lines[2])
	rw.MustWriteString(w, p.Lines[3])

	return rw.Result()
}

// ItemData is a server to client packet.
//
// Specifies complex data on an item; currently used only for maps.
//
// Maps If the first byte of the text is 0, the next two bytes are X start and Y
// start and the rest of the bytes are the colors in that column.
//
// If the first byte of the text is 1, the rest of the bytes are in groups of
// three: (data, x, y). The lower half of the data is the type (always 0 under
// vanilla) and the upper half is the direction.
//
// Total Size: 7 bytes + Text length
type ItemData struct {
	Type   int16
	ItemId int16  // Damage value
	Length int16  // len(Text)
	Text   []byte // ASCII text
}

func (p ItemData) Id() byte { return 0x83 }
func (p *ItemData) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Type = rw.MustReadShort(r)
	p.ItemId = rw.MustReadShort(r)
	p.Length = rw.MustReadShort(r)
	p.Text = rw.MustReadByteArray(r, int(p.Length))

	return rw.Result()
}

func (p *ItemData) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, p.Type)
	rw.MustWriteShort(w, p.ItemId)
	rw.MustWriteShort(w, p.Length)
	rw.MustWriteByteArray(w, p.Text)

	return rw.Result()
}

// TileEntityUpdate is a server to client packet.
// Total Size: 12 + itemstack bytes
type TileEntityUpdate struct {
	X      int32
	Y      int16
	Z      int32
	Action int8 // 1: set mob displayed inside mob spawner
	// Length int16 // Hidden in slice
	Data []byte // NBT Byte Array; iff Length > 0
}

func (p TileEntityUpdate) Id() byte { return 0x84 }
func (p *TileEntityUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.X = rw.MustReadInt(r)
	p.Y = rw.MustReadShort(r)
	p.Z = rw.MustReadInt(r)
	p.Action = rw.MustReadByte(r)
	length := rw.MustReadShort(r)
	p.Data = rw.MustReadByteArray(r, int(length))

	return rw.Result()
}

func (p *TileEntityUpdate) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.X)
	rw.MustWriteShort(w, p.Y)
	rw.MustWriteInt(w, p.Z)
	rw.MustWriteByte(w, p.Action)
	rw.MustWriteShort(w, int16(len(p.Data)))
	rw.MustWriteByteArray(w, p.Data)

	return rw.Result()
}

// StatIncrement is a server to client packet.
// Total Size: 6 bytes
type StatIncrement struct {
	StatId int32
	Amount int8
}

func (p StatIncrement) Id() byte { return 0xC8 }
func (p *StatIncrement) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.StatId = rw.MustReadInt(r)
	p.Amount = rw.MustReadByte(r)

	return rw.Result()
}

func (p *StatIncrement) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteInt(w, p.StatId)
	rw.MustWriteByte(w, p.Amount)

	return rw.Result()
}

// PlayerTabListPing is a server to client packet.
//
// NMS: sends one packet per user per tick (amounting to 20 packets/s for 1
// online user).
//
// Total Size: 6 bytes + length of string
type PlayerTabListPing struct {
	Name   string // Supports chat colouring, max 16 chars
	Online bool   // false = client will remove user from player list
	Ping   int16  // Ping, presumably in ms
}

func (p PlayerTabListPing) Id() byte { return 0xC9 }
func (p *PlayerTabListPing) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	p.Online = rw.MustReadBool(r)
	p.Ping = rw.MustReadShort(r)

	return rw.Result()
}

func (p *PlayerTabListPing) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteBool(w, p.Online)
	rw.MustWriteShort(w, p.Ping)

	return rw.Result()
}

// PlayerAbilities is a two-way packet.

// The latter 2 bytes are used to indicate the walking and flying speeds respectively, while the first byte is used to determine the value of 4 booleans.

// These booleans are whether damage is disabled (god mode, '8' bit), whether
// the player can fly ('4' bit), whether the player is flying ('2' bit), and
// whether the player is in creative mode ('1' bit).

// To get the values of these booleans, simply AND (&) the byte with 1,2,4 and 8
// respectively, to get the 0 or 1 bitwise value. To set them OR (|) them with
// their repspective masks. The vanilla client sends this packet when the player
// starts/stops flying with the second parameter changed accordingly. All other
// parameters are ignored by the vanilla server.

// Total Size: 4 bytes
type PlayerAbilities struct {
	GodMode, CanFly, Flying, Creative bool // Sent as byte
	FlyingSpeed                       int8
	WalkingSpeed                      int8
}

func (p PlayerAbilities) Id() byte { return 0xCA }
func (p *PlayerAbilities) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	flags := rw.MustReadByte(r)
	p.FlyingSpeed = rw.MustReadByte(r)
	p.WalkingSpeed = rw.MustReadByte(r)

	p.GodMode = flags&8 == 8
	p.CanFly = flags&4 == 4
	p.Flying = flags&2 == 2
	p.Creative = flags&1 == 1

	return rw.Result()
}

func (p *PlayerAbilities) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))

	var flags int8
	if p.GodMode {
		flags |= 8
	}
	if p.CanFly {
		flags |= 4
	}
	if p.Flying {
		flags |= 2
	}
	if p.Creative {
		flags |= 1
	}

	rw.MustWriteByte(w, flags)
	rw.MustWriteByte(w, p.FlyingSpeed)
	rw.MustWriteByte(w, p.WalkingSpeed)

	return rw.Result()
}

// TabComplete is a two-way packet.
// Total Size: 3 bytes + length of string
type TabComplete struct {
	Text string
}

func (p TabComplete) Id() byte { return 0xCB }
func (p *TabComplete) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Text = rw.MustReadString(r)

	return rw.Result()
}

func (p *TabComplete) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Text)

	return rw.Result()
}

// ClientSettings is a client to server packet.
// Total Size: 7 bytes + length of string
type ClientSettings struct {
	Locale       string
	ViewDistance int8 // 0-3 for "far", "normal", "short", "tiny"
	ChatFlags    int8 // Chat settings
	Difficulty   int8 // Client-side difficulty from options.txt
	ShowCape     bool // Client-side "show cape" option
}

func (p ClientSettings) Id() byte { return 0xCC }
func (p *ClientSettings) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Locale = rw.MustReadString(r)
	p.ViewDistance = rw.MustReadByte(r)
	p.ChatFlags = rw.MustReadByte(r)
	p.Difficulty = rw.MustReadByte(r)
	p.ShowCape = rw.MustReadBool(r)

	return rw.Result()
}

func (p *ClientSettings) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Locale)
	rw.MustWriteByte(w, p.ViewDistance)
	rw.MustWriteByte(w, p.ChatFlags)
	rw.MustWriteByte(w, p.Difficulty)
	rw.MustWriteBool(w, p.ShowCape)

	return rw.Result()
}

// ClientStatuses is a client to server packet.
// Total Size: 2 bytes
type ClientStatuses struct {
	Payload int8 // Bit field. 0: Initial spawn, 1: Respawn after death
}

func (p ClientStatuses) Id() byte { return 0xCD }
func (p *ClientStatuses) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Payload = rw.MustReadByte(r)

	return rw.Result()
}

func (p *ClientStatuses) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Payload)

	return rw.Result()
}

// ScoreObjective is a server to client packet.
// Total Size: 6 bytes + length of string
type ScoreObjective struct {
	Name   string // Unique name
	Value  string // The text to be displayed for the score
	Action int8   // 0 = create, 1 = remove, 2 = update display text
}

func (p ScoreObjective) Id() byte { return 0xCE }
func (p *ScoreObjective) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	p.Value = rw.MustReadString(r)
	p.Action = rw.MustReadByte(r)

	return rw.Result()
}

func (p *ScoreObjective) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteString(w, p.Value)
	rw.MustWriteByte(w, p.Action)

	return rw.Result()
}

// ScoreUpdate is a server to client packet.
// Total Size: 9 bytes + length of strings
type ScoreUpdate struct {
	EntityName string // Unique name to be displayed in the list.
	Action     int8   // 0 = create/update item, 1 = remove item
	Objetive   string // Unique name. Sent iff Action != 1
	Value      int32  // Sent iff Action != 1
}

func (p ScoreUpdate) Id() byte { return 0xCF }
func (p *ScoreUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.EntityName = rw.MustReadString(r)
	p.Action = rw.MustReadByte(r)
	p.Objetive = rw.MustReadString(r)
	p.Value = rw.MustReadInt(r)

	return rw.Result()
}

func (p *ScoreUpdate) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.EntityName)
	rw.MustWriteByte(w, p.Action)
	rw.MustWriteString(w, p.Objetive)
	rw.MustWriteInt(w, p.Value)

	return rw.Result()
}

// ScoreDisplay is a server to client packet.
// Total Size: 4 bytes + length of string
type ScoreDisplay struct {
	Position    int8   // 0 = list, 1 = sidebar, 2 = belowName
	DisplayName string // Unique name for the scoreboard to be displayed.
}

func (p ScoreDisplay) Id() byte { return 0xD0 }
func (p *ScoreDisplay) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Position = rw.MustReadByte(r)
	p.DisplayName = rw.MustReadString(r)

	return rw.Result()
}

func (p *ScoreDisplay) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Position)
	rw.MustWriteString(w, p.DisplayName)

	return rw.Result()
}

// ScoreTeams is a server to client packet.
// Total Size: Variable
type ScoreTeams struct {
	Name                        string   // Unique, shared with scoreboard
	Mode                        int8     // 0: create, 1: remove, 2: update, 3: player join, 4: player leave
	DisplayName, Prefix, Suffix string   // iff Mode == (0 | 2)
	FriendlyFire                int8     // iff Mode == (0 | 2); 0: off, 1: on, 3: see friendly invisibles
	Count                       int16    // iff Mode == (0 | 3 | 4)
	Players                     []string // iff Mode == (0 | 3 | 4)
}

func (p ScoreTeams) Id() byte { return 0xD1 }
func (p *ScoreTeams) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	p.Mode = rw.MustReadByte(r)
	p.DisplayName = rw.MustReadString(r)
	p.Prefix = rw.MustReadString(r)
	p.Suffix = rw.MustReadString(r)
	p.FriendlyFire = rw.MustReadByte(r)
	p.Count = rw.MustReadShort(r)
	p.Players = make([]string, p.Count)
	for i := 0; i < len(p.Players); i++ {
		p.Players[i] = rw.MustReadString(r)
	}

	return rw.Result()
}

func (p *ScoreTeams) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteByte(w, p.Mode)
	rw.MustWriteString(w, p.DisplayName)
	rw.MustWriteString(w, p.Prefix)
	rw.MustWriteString(w, p.Suffix)
	rw.MustWriteByte(w, p.FriendlyFire)
	rw.MustWriteShort(w, p.Count)
	for i := 0; i < len(p.Players); i++ {
		rw.MustWriteString(w, p.Players[i])
	}

	return rw.Result()
}

// PluginMessage is a two-way packet.
// Total Size: 5 bytes + len(Name) + len(Payload)
type PluginMessage struct {
	Name string
	// Length  int16 // len(Payload)
	Payload []byte
}

func (p PluginMessage) Id() byte { return 0xFA }
func (p *PluginMessage) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Name = rw.MustReadString(r)
	length := int(rw.MustReadShort(r))
	p.Payload = rw.MustReadByteArray(r, length)

	return rw.Result()
}

func (p *PluginMessage) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Name)
	rw.MustWriteShort(w, int16(len(p.Payload)))
	rw.MustWriteByteArray(w, p.Payload)

	return rw.Result()
}

// EncryptionKeyResponse is a two-way packet.
// Total Size: 5 bytes + len(Secret) + len(Token)
type EncryptionKeyResponse struct {
	// SecretLength, Tokenlength int16 // stored in slice len fn
	Secret, Token []byte
}

func (p EncryptionKeyResponse) Id() byte { return 0xFC }
func (p *EncryptionKeyResponse) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	secretLen := rw.MustReadShort(r)
	p.Secret = rw.MustReadByteArray(r, int(secretLen))
	tokenLen := rw.MustReadShort(r)
	p.Token = rw.MustReadByteArray(r, int(tokenLen))

	return rw.Result()
}

func (p *EncryptionKeyResponse) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteShort(w, int16(len(p.Secret)))
	rw.MustWriteByteArray(w, p.Secret)
	rw.MustWriteShort(w, int16(len(p.Token)))
	rw.MustWriteByteArray(w, p.Token)

	return rw.Result()
}

// EncryptionKeyRequest is a server to client packet.
// Total Size: 7 bytes + length of string + length of key + length of token
type EncryptionKeyRequest struct {
	ServerId string
	// SecretLength, Tokenlength int16 // stored in slice len fn
	PublicKey, Token []byte
}

func (p EncryptionKeyRequest) Id() byte { return 0xFD }
func (p *EncryptionKeyRequest) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.ServerId = rw.MustReadString(r)
	pkLen := int(rw.MustReadShort(r))
	p.PublicKey = rw.MustReadByteArray(r, pkLen)
	tokenLen := int(rw.MustReadShort(r))
	p.Token = rw.MustReadByteArray(r, tokenLen)

	return rw.Result()
}

func (p *EncryptionKeyRequest) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.ServerId)
	rw.MustWriteShort(w, int16(len(p.PublicKey)))
	rw.MustWriteByteArray(w, p.PublicKey)
	rw.MustWriteShort(w, int16(len(p.Token)))
	rw.MustWriteByteArray(w, p.Token)

	return rw.Result()
}

// ServerListPing is a client to server packet.
// Total Size: 2 bytes
type ServerListPing struct {
	Magic int8 // always 1
}

func (p ServerListPing) Id() byte { return 0xFE }
func (p *ServerListPing) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Magic = rw.MustReadByte(r)

	return rw.Result()
}

func (p *ServerListPing) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteByte(w, p.Magic)

	return rw.Result()
}

// Disconnect is a two-way packet.
// Total Size: 3 bytes + length of strings
type Disconnect struct {
	Reason string // Displayed to the client when the connection terminates
}

func (p Disconnect) Id() byte { return 0xFF }
func (p *Disconnect) ReadFrom(r io.Reader) (n int64, err error) {
	rw.Reset()
	// id := rw.MustReadByte(r)
	// rw.Check(CheckPacketId(id))
	p.Reason = rw.MustReadString(r)

	return rw.Result()
}

func (p *Disconnect) WriteTo(w io.Writer) (n int64, err error) {
	rw.Reset()
	id := Id(p.Id())
	rw.Must(id.WriteTo(w))
	rw.MustWriteString(w, p.Reason)

	return rw.Result()
}

type Id byte

func (i *Id) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, i)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (i *Id) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, *i)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
