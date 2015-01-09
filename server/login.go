package server

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/minero/minero/constants/biome"
	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
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
	r = &packet.EntityNamedSpawn{
		Entity: sender.Id(),
		Name:   sender.Name,
		X:      startX, Y: startY, Z: startZ,
		Yaw:      0.0,
		Pitch:    0.0,
		Item:     0,
		Metadata: player.JustLoginMetadata(sender.Name),
	}
	s.BroadcastPacket(r)

	// Instantiate all other users on new client
	s.BroadcastLogin(sender)

	// Save player to server list
	s.AddPlayer(sender)
}

// type Chunk struct {
// 	Block    [16][16][16][16]byte // 16 sections of 4k each. Order YZX.
// 	Meta     [16][16][16][8]byte  // 16 sections of 2k each. Order YZX.
// 	Light    [16][16][16][8]byte  // 16 sections of 2k each. Order YZX.
// 	SkyLight [16][16][16][8]byte  // 16 sections of 2k each. Order YZX.
// 	Add      [16][16][16][8]byte  // 16 sections of 2k each. Order YZX.
// 	Biome    [16][16]byte         // 16 x 16 bytes. Order XZ.
// }

func VirtualChunks(x, z, height int32) packet.Packet {
	var temp bytes.Buffer
	var nul [1 << 11]byte

	// Block type. 16*4k bytes. Order YZX. Stone.
	for i := 0; i < sections(height); i++ {
		temp.Write(bytes.Repeat([]byte{1}, 1<<12)) // 4x4k sections of stone.
	}

	// Block metadata. 16*2k bytes. Zeros.
	for i := 0; i < sections(height); i++ {
		temp.Write(nul[:])
	}
	// Block light. 16*2k bytes. Zeros.
	for i := 0; i < sections(height); i++ {
		temp.Write(nul[:])
	}
	// Block sklylight. 16*2k bytes.
	for i := 0; i < sections(height); i++ {
		temp.Write(bytes.Repeat([]byte{0xcc}, 1<<11)) // Light all over the place.
	}
	// Add mask. 16*2k bytes. Zeros. Only sent iff Add != 0
	// for i := 0; i < sections(height); i++ {
	// 	temp.Write(nul[:])
	// }

	// Biomes. 256 bytes. Plains.
	temp.Write(bytes.Repeat([]byte{biome.Plains}, 1<<8))

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	io.Copy(zw, &temp)
	zw.Close()

	return &packet.ChunkData{
		X:              x,
		Z:              z,
		AllColSections: true,
		Primary:        sectionMask(height),
		Add:            0,
		ChunkData:      buf.Bytes(),
	}
}

func sections(height int32) int {
	s := int(height / 16)
	if s == 0 {
		return 1
	}
	return s
}

func sectionMask(height int32) uint16 {
	return uint16(1<<uint(sections(height)) - 1)
}
