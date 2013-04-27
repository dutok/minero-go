package server

import (
	"log"

	"github.com/toqueteos/minero/material"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0F handles incoming requests of packet 0x0F: PlayerBlockPlace
func Handle0F(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerBlockPlace)
	pkt.ReadFrom(sender.Conn)

	log.Printf("PlayerBlockPlace: %+v", pkt)

	var bid int16
	if pkt.HeldItem != nil {
		bid = pkt.HeldItem.BlockId
		if pkt.HeldItem.BlockId == -1 {
			bid = int16(material.Stone)
		}
	}

	r := &packet.BlockChange{
		X:         pkt.X, // int32
		Y:         pkt.Y, // int8
		Z:         pkt.Z, // int32
		BlockType: bid,   // int16 // New block type for block
		BlockMeta: 0,     // int8
	}
	server.BroadcastPacket(r)
}
