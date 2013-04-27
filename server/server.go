package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"sync"
	"time"

	"github.com/toqueteos/minero/command"
	"github.com/toqueteos/minero/config"
	"github.com/toqueteos/minero/proto/auth"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/list/players"
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

	cmdList map[string]command.Cmder

	// Embed list handlers
	players.PlayersList
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
		id:      serverId(),
		config:  c,
		privKey: auth.GenerateKeyPair(),

		// Load from config
		Motd: c.Get("server.motd"),

		PlayersList: players.New(),
	}

	return s
}

// Id returns server's Id.
func (s Server) Id() string { return s.id }

func (s Server) PlayersIn() int { return s.PlayersList.Len() }

// PublicKey returns the ASN.1 encoded version of server's x.509 public key.
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

	for !s.working {
		// Wait for a connection.
		conn, err := s.net.Accept()
		if err != nil {
			log.Println(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handle(conn)
	}
}

func (s *Server) handle(c net.Conn) {
	defer c.Close()
	defer log.Println("Connection closed:", c.RemoteAddr())
	log.Println("Got connection from:", c.RemoteAddr())

	// Create player "instance" and save it to player list
	p := player.New(c)
	s.AddPlayer(p)

	// Ensure player is deleted from online list, doesn't care about why he/she
	// disconnects.
	defer s.RemPlayer(p)

	// Send KeepAlive packet every 30s (x20 in-game ticks)
	go func() {
		for _ = range time.Tick(30 * time.Second) {
			r := &packet.KeepAlive{RandomId: mrand.Int31()}
			r.WriteTo(p.Conn)
		}
	}()

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
	return fmt.Sprintf("minero%x-%d", mrand.Int31(), time.Now().Year())
}
