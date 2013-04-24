package player

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/toqueteos/minero/id"
	"github.com/toqueteos/minero/proto/packet"
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

func (p *Player) SetPos(x, y, z float64) {
	p.Lock()
	p.X = x
	p.Y = y
	p.Z = z
	p.Unlock()
}
func (p *Player) SetLook(pitch, yaw float32) {
	p.Lock()
	p.Pitch = pitch
	p.Yaw = yaw
	p.Unlock()
}

func (p *Player) SetReady() {
	p.Lock()
	p.Ready = true
	p.Unlock()
}

func (p *Player) SendMessage(msg string) {
	pkt := packet.ChatMessage{msg}
	pkt.WriteTo(p.Conn)
}

func (p *Player) OnlineMode(m bool, secret []byte) {
	p.crypto = m
	if m {
		p.Conn = cfb8.New(p.net, secret)
	}
}
