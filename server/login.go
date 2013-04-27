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

func (s *Server) HandleLogin(sender *player.Player) {
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

	// BUG(toqueteos): Load nearby chunks
	for z := int32(-1); z < 2; z++ {
		for x := int32(-1); x < 2; x++ {
			VirtualChunks(x, z, 64).WriteTo(sender.Conn)
		}
	}

	const (
		startX = 8.0
		startY = 65.0
		startZ = 8.0
	)

	// Client's spawn position
	r = &packet.SpawnPosition{X: int32(startX), Y: int32(startY), Z: int32(startZ)}
	r.WriteTo(sender.Conn)

	// Client Pos & Look
	r = &packet.PlayerPosLook{
		startX, startY, startZ, // X, Y, Z
		startY + 1.6, // Stance
		0.0, 0.0,     // Yaw + Pitch
		true, // OnGround
	}
	r.WriteTo(sender.Conn)

	// Send nearby clients new client's info
	meta := mct.NewMetadata()
	meta.Entries[0] = &mct.EntryByte{0}
	_ = meta
	r = &packet.EntityNamedSpawn{
		Entity: sender.Id(),
		Name:   sender.Name,
		X:      startX, Y: startY, Z: startZ,
		Yaw:      0.0,
		Pitch:    0.0,
		Item:     0,
		Metadata: mct.MetadataFrom([]byte{0, 0, 6, 0, 127}),
	}
	s.BroadcastPacket(r)

	// Initialize entity on other player's clients
	// r = &packet.Entity{sender.Id()}
	// server.BroadcastPacket(r)

	// Send nearby clients client's Pos & Look
	// r = &packet.EntityTeleport{
	// 	Entity: sender.Id(),
	// 	X:      startX, Y: startY, Z: startZ,
	// }
	// s.BroadcastPacket(r)

	// r = &packet.EntityHeadLook{
	// 	Entity:  sender.Id(),
	// 	HeadYaw: 0.0,
	// }
	// s.BroadcastPacket(r)
}

func VirtualChunks(x, z, height int32) packet.Packet {
	var temp bytes.Buffer
	var nul [1 << 11]byte

	temp.Write(bytes.Repeat([]byte{1}, 1<<12))  // Block type. 4k bytes. Order YZX. Stone 1x
	temp.Write(nul[:])                          // Block metadata. 2k bytes. Zeros
	temp.Write(nul[:])                          // Block light. 2k bytes. Zeros
	temp.Write(bytes.Repeat([]byte{12}, 1<<12)) // Block sklylight. 2k bytes.
	temp.Write(nul[:])                          // Add mask. 2k bytes. Zeros

	log.Println("AFTER Repeat:", temp.Len())

	// Repeat for each section
	temp.Write(bytes.Repeat(temp.Bytes(), int(sections(height)-1)))

	log.Println("BEFORE Repeat:", temp.Len())

	temp.Write(bytes.Repeat([]byte{1}, 1<<8)) // Biomes. 256 bytes. Plains

	log.Println("BEFORE Repeat + Biomes:", temp.Len())

	log.Println("Sections:", sections(height))

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	io.Copy(zw, &temp)
	zw.Close()

	log.Println("ZLIB chunk size:", buf.Len())

	return &packet.ChunkData{
		X:              x,
		Z:              z,
		AllColSections: true,
		Primary:        (1 << sections(height)) - 1,
		Add:            0,
		ChunkData:      buf.Bytes(),
	}
}

func sections(height int32) uint {
	s := uint(height / 16)
	if s == 0 {
		return 1
	}
	return s
}
