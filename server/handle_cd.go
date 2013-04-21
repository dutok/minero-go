package server

import (
	"log"
	"time"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCD handles incoming requests of packet 0xCD: ClientStatuses
func HandleCD(s *Server, player *player.Player) {
	p := new(packet.ClientStatuses)
	p.ReadFrom(player.Conn)

	switch p.Payload {
	case 0:
		var r packet.Packet

		r = &packet.LoginInfo{
			Entity:     33,
			LevelType:  "default",
			GameMode:   1,
			Dimension:  0,
			Difficulty: 2,
			MaxPlayers: 32,
		}
		r.WriteTo(player.Conn)

		r = &packet.SpawnPosition{X: 0, Y: 64, Z: 0}
		r.WriteTo(player.Conn)

		r = &packet.PlayerPosLook{
			X:        0.0,
			Y:        64.0,
			Z:        0.0,
			Stance:   65.6,
			Yaw:      0.0,
			Pitch:    0.0,
			OnGround: true,
		}
		r.WriteTo(player.Conn)

		// Save player login timestamp
		player.Since = time.Now().Unix()
	case 1:
	default:
		log.Println("Weird packet 0xCB payload:", p.Payload)
		r := packet.Disconnect{"Weird packet 0xCB payload"}
		r.WriteTo(player.Conn)
	}
}
