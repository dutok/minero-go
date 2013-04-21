package server

import (
	"log"

	"github.com/toqueteos/minero"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle02 handles incoming requests of packet 0x02: Handshake
func Handle02(s *Server, player *player.Player) {
	p := new(packet.Handshake)
	p.ReadFrom(player.Conn)

	log.Printf("Handshake from: %q [%s]", p.Username, player.Net.RemoteAddr())

	if p.Version != minero.ProtoNum {
		log.Printf("Wrong Protocol version. Player: %d, Server: %d\n",
			p.Version, minero.ProtoNum)
		return
	}
	player.Name = p.Username

	// BUG(toqueteos): OnlineMode = false sends 0x01 packet here
	if false {
		// ...
	} else {
		// Succesful handshake, prepare Encryption Request
		r := packet.EncryptionKeyRequest{
			ServerId:  s.Id(),
			PublicKey: s.PublicKey(),
			Token:     s.Token(),
		}
		r.WriteTo(player.Conn)
		player.Token = r.Token
	}
}
