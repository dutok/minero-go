package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCB handles incoming requests of packet 0xCB: TabComplete
func HandleCB(server *Server, sender *player.Player) {
	pkt := new(packet.TabComplete)
	pkt.ReadFrom(sender.Conn)

	log.Printf("TabComplete: %+v", pkt)
}
