package player

import (
	"io"
	"net"

	"github.com/toqueteos/minero/util/crypto/cfb8"
)

type Player struct {
	Net  net.Conn
	Conn io.ReadWriter

	Secret, Token []byte
	Crypto        bool

	Name  string
	Since int64
}

func New(c net.Conn) *Player      { return &Player{Net: c, Conn: c} }
func (p Player) String() string   { return p.Name }
func (p Player) OnlineFor() int64 { return p.Since } // Duration?

func (p *Player) OnlineMode(mode bool) {
	p.Crypto = mode
	if mode {
		p.Conn = cfb8.New(p.Net, p.Secret)
	}
}
