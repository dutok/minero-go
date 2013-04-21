package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	prand "math/rand"
	"net"
	"time"

	"github.com/toqueteos/minero/proto/auth"
	"github.com/toqueteos/minero/server/player"
)

type Server struct {
	privKey *rsa.PrivateKey
	pubEx   []byte
	token   []byte

	id         string
	host, port string
	motd       string

	players map[string]*player.Player
	in, max int // max players
}

func New(host, port string) *Server {
	return &Server{
		privKey: auth.GenerateKeyPair(),

		id:   serverId(),
		host: host,
		port: port,
		motd: "Minero server", // Read this from config

		players: make(map[string]*player.Player),
		max:     64, // Read this from config
	}
}

func (s Server) Id() string    { return s.id }
func (s Server) Motd() string  { return s.motd }
func (s Server) Token() []byte { return auth.EncryptionBytes() }
func (s Server) Host() string  { return s.host }
func (s Server) Port() string  { return s.port }

func (s *Server) PublicKey() []byte {
	if s.pubEx == nil {
		var err error
		s.pubEx = auth.KeyExchange(&s.privKey.PublicKey)
		if s.pubEx == nil {
			log.Fatal("Couldn't marshal public key:", err)
			return nil
		}
	}
	return s.pubEx
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
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
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

	// var buf = bufio.NewReader(c)
	var buf = make([]byte, 1)

	// Create player "instance"
	p := player.New(c)

	// Save it after successful login
	// s.players[c.RemoteAddr().String()] = p
	// s.in++
	// defer func() {
	// 	delete(s.players, c.RemoteAddr().String())
	// 	s.in--
	// }()

	for {
		n, err := p.Conn.Read(buf)
		if n != 1 || err != nil {
			log.Println("Server.Handle.ReadByte:", err)
			return
		}
		pid := buf[0]
		log.Printf("Packet: %#x\n", pid)

		h := HandlerFor[pid]
		if h != nil {
			h(s, p)
		}
	}
}

// func (s *Server) AddPlayer(c net.Conn) {
// 	s.players[c.RemoteAddr().String()] = new(player.Player)
// }

// func (s *Server) GetPlayer(c net.Conn) *player.Player {
// 	return s.players[c.RemoteAddr().String()]
// }

func serverId() string {
	return fmt.Sprintf("minero%x-%d", prand.Int31(), time.Now().Year())
}

func yesno(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
