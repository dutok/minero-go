package player

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/toqueteos/minero/id"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/util"
	"github.com/toqueteos/minero/util/crypto/cfb8"
)

type Player struct {
	sync.Mutex
	net  net.Conn
	Conn io.ReadWriter

	Token  []byte
	crypto bool

	Ready bool

	Name     string
	since    int64
	eid      int32
	GameMode int8

	X, Y, Z    float64
	Pitch, Yaw float32
}

func New(c net.Conn) *Player {
	return &Player{
		net:   c,
		Conn:  c,
		since: time.Now().Unix(),
		eid:   id.Get(),
	}
}

func (p Player) String() string     { return fmt.Sprintf("Player{%q#%d}", p.Name, p.eid) }
func (p Player) RemoteAddr() string { return p.net.RemoteAddr().String() }
func (p Player) Id() int32          { return p.eid }
func (p Player) OnlineSince() int64 { return p.since }
func (p Player) UsesCrypto() bool   { return p.crypto }

// OnlineMode set's authentication on for p. Authentication is required for
// online_mode=true servers.
func (p *Player) OnlineMode(m bool, secret []byte) {
	p.crypto = m
	if m {
		p.Conn = cfb8.New(p.net, secret)
	}
}

func (p *Player) SetReady() { p.Lock(); p.Ready = true; p.Unlock() }

func (p *Player) SetPos(x, y, z float64) {
	p.Lock()
	p.X = x
	p.Y = y
	p.Z = z
	p.Unlock()
}

func (p *Player) SetX(x float64) { p.Lock(); p.X = x; p.Unlock() }
func (p *Player) SetY(y float64) { p.Lock(); p.Y = y; p.Unlock() }
func (p *Player) SetZ(z float64) { p.Lock(); p.Z = z; p.Unlock() }

func (p *Player) SetLook(pitch, yaw float32) {
	p.Lock()
	p.Pitch = pitch
	p.Yaw = yaw
	p.Unlock()
}

func (p *Player) SetPitch(pitch float32) { p.Lock(); p.Pitch = pitch; p.Unlock() }
func (p *Player) SetYaw(yaw float32)     { p.Lock(); p.Yaw = yaw; p.Unlock() }

// SendMessage sends a chat message to p.
func (p *Player) SendMessage(message string) {
	pkt := &packet.ChatMessage{message}
	pkt.WriteTo(p.Conn)
}

// BroadcastMessage sends a message to all `ready` players from toList. If p is
// in that list he/she is ommited.
func (p *Player) BroadcastMessage(toList map[string]*Player, message string) {
	for _, to := range toList {
		// Send message only to other ready players
		if to.Ready && to.Name != p.Name {
			to.SendMessage(message)
		}
	}
}

// Tick sends a KeepAlive packet every 1000 in-game ticks (50s).
func (p *Player) Tick(t int64) {
	if t%util.Ticks(50) == 0 {
		r := &packet.KeepAlive{RandomId: rand.Int31()}
		r.WriteTo(p.Conn)
	}
}

// Destroy attempts to release all resources allocated by this player.
func (p *Player) Destroy() {
	id.Rel(p.Id())
}
