package server

import (
	"fmt"
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleFF handles incoming requests of packet 0xFF: Disconnect
func HandleFF(server *Server, sender *player.Player) {
	pkt := new(packet.Disconnect)
	pkt.ReadFrom(sender.Conn)

	log.Printf("Player %q exit. Reason: %s", sender.Name, pkt.Reason)

	// Send message to all other players
	msg := fmt.Sprintf("%s disconnected.", sender.Name)
	server.BroadcastMessage(msg)
}
