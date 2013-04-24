package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	prand "math/rand"
	"net"
	"sync"
	"time"

	"github.com/toqueteos/minero/command"
	"github.com/toqueteos/minero/config"
	"github.com/toqueteos/minero/proto/auth"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

type Server struct {
	sync.Mutex

	id string

	net     net.Listener
	working bool

	config *config.Config

	privKey *rsa.PrivateKey
	pubKey  []byte
	token   []byte

	// Message of the day. Text appears on server list.
	Motd string

	// Stop message. Text appears on server list.
	Stop string

	cmdList    map[string]command.Cmder
	playerList map[string]*player.Player
	// pluginList map[string]*plugin.Plugin
	// worldList  map[string]*world.World
}

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
		id:      serverId(),
		config:  c,
		privKey: auth.GenerateKeyPair(),

		playerList: make(map[string]*player.Player),
	}

	return s
}

func (s Server) Id() string    { return s.id }
func (s Server) Token() []byte { return auth.EncryptionBytes() }
func (s Server) CmdManager()   {}

func (s *Server) PublicKey() []byte {
	if s.pubKey == nil {
		var err error
		s.pubKey = auth.KeyExchange(&s.privKey.PublicKey)
		if s.pubKey == nil {
			log.Fatal("Couldn't marshal public key:", err)
			return nil
		}
	}
	return s.pubKey
}

// Decrypt decrypts whatever the client encrypted with its keypair.
func (s *Server) Decrypt(what []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, s.privKey, what)
}

func (s *Server) CheckUser(name string, secret []byte) bool {
	r, err := auth.CheckUser(name, s.id, secret, s.PublicKey())
	if err != nil {
		return false
	}
	return r
}

func (s *Server) Run() {
	var err error

	addr := s.config.Get("server.host") + ":" + s.config.Get("server.port")
	log.Printf("Listening on address: %q", addr)

	s.net, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer s.net.Close()

	for !s.working {
		// Wait for a connection.
		conn, err := s.net.Accept()
		if err != nil {
			log.Println(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.Handle(conn)
	}
}

func (s *Server) Handle(c net.Conn) {
	defer c.Close()
	defer log.Println("Connection closed:", c.RemoteAddr())
	log.Println("Got connection from:", c.RemoteAddr())

	// Create player "instance"
	p := player.New(c)

	// Send KeepAlive packet every 35s (x20 in-game ticks)
	go func() {
		for _ = range time.Tick(35 * time.Second) {
			r := &packet.KeepAlive{RandomId: prand.Int31()}
			r.WriteTo(p.Conn)
		}
	}()

	var buf = make([]byte, 1)
	for {
		n, err := p.Conn.Read(buf)
		if n != 1 || err != nil {
			// log.Println("Server.Handle.ReadByte:", err)
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

func (s *Server) BroadcastOthers(p *player.Player, pkt packet.Packet) {
	for _, pl := range s.playerList {
		if pl.Ready && p.Name != pl.Name {
			pkt.WriteTo(pl.Conn)
		}
	}
}
func (s *Server) BroadcastMessage(msg string) {
	for _, p := range s.playerList {
		if p.Ready {
			p.SendMessage(msg)
		}
	}
}

// AddPlayer
func (s *Server) AddPlayer(p *player.Player) {
	s.Lock()
	s.playerList[p.Name] = p
	s.Unlock()
}

// RemPlayer
func (s *Server) RemPlayer(p *player.Player) {
	s.Lock()
	delete(s.playerList, p.Name)
	s.Unlock()
}

func serverId() string {
	return fmt.Sprintf("minero%x-%d", prand.Int31(), time.Now().Year())
}

func yesno(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
