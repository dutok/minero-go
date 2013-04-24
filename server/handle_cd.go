package server

import (
	"bytes"
	"compress/zlib"
	"io"
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
	mct "github.com/toqueteos/minero/types/minecraft"
)

// HandleCD handles incoming requests of packet 0xCD: ClientStatuses
func HandleCD(server *Server, sender *player.Player) {
	pkt := new(packet.ClientStatuses)
	pkt.ReadFrom(sender.Conn)

	switch pkt.Payload {
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
		r.WriteTo(sender.Conn)

		// -----

		var temp bytes.Buffer
		var nul [1 << 12]byte

		// Block type. 4k bytes. Order YZX
		temp.Write(bytes.Repeat([]byte{7}, 1<<8))  // Bedrock 1x1x1
		temp.Write(bytes.Repeat([]byte{1}, 1<<12)) // Stone 1x
		// Block metadata. 2k bytes.
		temp.Write(nul[:1<<11]) // zeros
		// Block light. 2k bytes.
		temp.Write(nul[:1<<11]) // zeros
		// Block sklylight. 2k bytes.
		temp.Write(bytes.Repeat([]byte{12}, 128))
		temp.Write(nul[:1<<11-1<<7]) // zeros
		// Add mask. 2k bytes.
		// temp.Write(nul[:1<<12]) // zeros

		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		io.Copy(zw, &temp)
		zw.Close()

		ch := &packet.ChunkData{
			X:              0,
			Z:              0,
			AllColSections: true,
			Primary:        1,
			Add:            0,
			ChunkData:      buf.Bytes(),
		}
		ch.WriteTo(sender.Conn)

		// -----

		r = &packet.SpawnPosition{X: 8, Y: 18, Z: 8}
		r.WriteTo(sender.Conn)

		r = &packet.PlayerPosLook{
			X:        8.0,
			Y:        18.0,
			Z:        8.0,
			Stance:   19.6,
			Yaw:      0.0,
			Pitch:    0.0,
			OnGround: true,
		}
		r.WriteTo(sender.Conn)

		meta := mct.NewMetadata()
		meta.Entries[0] = &mct.EntryByte{0}
		r = &packet.EntityNamedSpawn{
			Entity:   sender.Id(),
			Name:     sender.Name,
			X:        8.0,
			Y:        16.0,
			Z:        8.0,
			Yaw:      0.0,
			Pitch:    0.0,
			Item:     0,
			Metadata: meta,
		}
		server.BroadcastOthers(sender, r)
	case 1:
	default:
		log.Println("Weird packet 0xCB payload:", pkt.Payload)
		r := packet.Disconnect{"Weird packet 0xCB payload"}
		r.WriteTo(sender.Conn)
	}
}
