package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// HandleCD handles incoming requests of packet 0xCD: ClientStatuses
func HandleCD(server *Server, sender *player.Player) {
	pkt := new(packet.ClientStatuses)
	pkt.ReadFrom(sender.Conn)

	switch pkt.Payload {
	case 0:
		server.HandleLogin(sender)
	case 1:
	default:
		log.Println("Weird packet 0xCB payload:", pkt.Payload)
		r := packet.Disconnect{"Weird packet 0xCB payload"}
		r.WriteTo(sender.Conn)
	}
}
