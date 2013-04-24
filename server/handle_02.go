package server

import (
	"log"

	"github.com/toqueteos/minero"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle02 handles incoming requests of packet 0x02: Handshake
func Handle02(server *Server, sender *player.Player) {
	pkt := new(packet.Handshake)
	pkt.ReadFrom(sender.Conn)

	log.Printf("Handshake from: %q [%s]", pkt.Username, sender.RemoteAddr())

	if pkt.Version != minero.ProtoNum {
		log.Printf("Wrong Protocol version. Player: %d, Server: %d\n",
			pkt.Version, minero.ProtoNum)
		return
	}

	// Get this player his own entity Id
	sender.NewEntityId(pkt.Username)
	// Save player to list
	server.AddPlayer(sender)

	// BUG(toqueteos): OnlineMode = false sends 0x01 packet here
	if false {
		// ...
	} else {
		// Succesful handshake, prepare Encryption Request
		r := packet.EncryptionKeyRequest{
			ServerId:  server.Id(),
			PublicKey: server.PublicKey(),
			Token:     server.Token(),
		}
		r.WriteTo(sender.Conn)
		sender.Token = r.Token
	}
}
