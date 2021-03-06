package server

import (
	"crypto/rsa"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/minero/minero/cmd"
	"github.com/minero/minero/config"
	"github.com/minero/minero/id"
	"github.com/minero/minero/proto/auth"
	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/list/players"
	"github.com/minero/minero/server/list/tickers"
	"github.com/minero/minero/server/player"
)

type Server struct {
	sync.Mutex

	id string

	net     net.Listener
	working bool

	ticks int64

	config *config.Config

	privateKey *rsa.PrivateKey
	publicKey  []byte
	token      []byte

	gameTime int64

	// Message of the day. Text appears on server list.
	Motd string
	// Stop message. Text appears on server list.
	Stop string

	// Embed list handlers
	Cmds map[string]cmd.Cmder
	players.Players
	tickers.Tickers
}

// New initializes a new server instance and loads server.conf file if one
// exists, otherwise it'll create a new one.
func New(c *config.Config) *Server {
	log.Println("Generating keypair.")

	// Generate config
	if c == nil {
		c = config.New()
		if err := c.ParseFile("./server.conf"); err != nil {
			c = ConfigCreate()
		}
	}

	s := &Server{
		id:         serverId(),
		config:     c,
		privateKey: auth.GenerateKeyPair(),

		// Load from config
		Motd: c.Get("server.motd"),

		Cmds:    make(map[string]cmd.Cmder),
		Players: players.New(),
		Tickers: tickers.New(),
	}

	return s
}

// Id returns server's Id.
func (s Server) Id() string { return s.id }

// PublicKey returns the ASN.1 encoded version of server's x.509 public key.
func (s *Server) PublicKey() []byte {
	if s.publicKey == nil {
		var err error
		s.publicKey = auth.KeyExchange(&s.privateKey.PublicKey)
		if s.publicKey == nil {
			log.Fatal("Couldn't marshal public key:", err)
			return nil
		}
	}
	return s.publicKey
}

// PrivateKey returns server's private key.
func (s *Server) PrivateKey() *rsa.PrivateKey { return s.privateKey }

// CheckUser check's if user is premium, only used when config var
// "server.online_mode" = true.
func (s *Server) CheckUser(name string, secret []byte) bool {
	r, err := auth.CheckUser(name, s.id, secret, s.PublicKey())
	if err != nil {
		return false
	}
	return r
}

// Run starts up the server.
func (s *Server) Run() {
	var err error

	addr := s.config.Get("server.host") + ":" + s.config.Get("server.port")
	log.Printf("Listening on address: %q", addr)

	s.net, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer s.net.Close()

	// Start server ticker
	go s.doTick()

	for !s.working {
		// Wait for a connection.
		conn, err := s.net.Accept()
		if err != nil {
			log.Println(err)
		}

		// Handle the connection in a new goroutine.
		// Multiple connections may be served concurrently.
		go s.handle(conn)
	}
}

// Tick updates time by one tick. Implements tick.Ticker interface.
func (s *Server) Tick(t int64) {
	pkt := &packet.TimeUpdate{
		Time:     s.gameTime,
		WorldAge: t,
	}
	s.BroadcastPacket(pkt)
}

func (s *Server) doTick() {
	// This id *NEVER* gets released
	s.AddTicker(id.Get(), s)
	for _ = range time.Tick(50 * time.Millisecond) {
		s.TickAll(s.ticks)
		s.ticks++
	}
}

func (s *Server) handle(c net.Conn) {
	defer c.Close()
	defer log.Println("Connection closed:", c.RemoteAddr())
	log.Println("Got connection from:", c.RemoteAddr())

	// Create player "instance", don't save it yet, it may be a ServerList Ping.
	p := player.New(c)
	s.AddTicker(p.Id(), p)

	// Ensure player is deleted from online list, doesn't care about why he/she
	// disconnects.
	defer p.Destroy()
	defer s.RemPlayer(p)
	defer s.RemTicker(p.Id())

	var buf = make([]byte, 1)
	for {
		n, err := p.Conn.Read(buf)
		if n != 1 || err != nil {
			return
		}
		pid := buf[0]

		h := HandlerFor[pid]
		if h != nil {
			h(s, p)
		} else {
			log.Fatalf("Can't handle packet %#x. Closing", pid)
		}
	}
}

// Kick kicks a player from the server
func (s *Server) Kick(p *player.Player) {
	p.SendMessage("You were kicked from the server.")
	msg := fmt.Sprintf("Player %q was kicked from the server.", p.Name)
	s.BroadcastMessage(msg)
}

func serverId() string {
	return fmt.Sprintf("minero%x-%d", rand.Int31(), time.Now().Year())
}
