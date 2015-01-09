package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// HandleCA handles incoming requests of packet 0xCA: PlayerAbilities
func HandleCA(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerAbilities)
	pkt.ReadFrom(sender.Conn)

	log.Printf("PlayerAbilities: %+v", pkt)
}
