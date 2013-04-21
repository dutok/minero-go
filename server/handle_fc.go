package server

import (
	"bytes"
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleFC handles incoming requests of packet 0xFC: EncryptionKeyResponse
func HandleFC(s *Server, player *player.Player) {
	p := new(packet.EncryptionKeyResponse)
	p.ReadFrom(player.Conn)

	// Decrypt shared secret and token with server's private key.
	var secret, token []byte
	// var err error
	secret, _ = s.Decrypt(p.Secret)
	token, _ = s.Decrypt(p.Token)

	// Ensure token matches
	if !bytes.Equal(token, player.Token) {
		log.Println("Tokens don't match.")
		r := &packet.Disconnect{Reason: ReasonPiratedGame}
		r.WriteTo(player.Conn)
		return
	}

	log.Println("Shared secret length:", len(secret))

	// Ensure player is legit
	if !s.CheckUser(player.Name, secret) {
		log.Println("Failed to verify username!")
		r := packet.Disconnect{"Failed to verify username!"}
		r.WriteTo(player.Conn)
		return
	}

	// Send empty EncryptionKeyResponse
	r := new(packet.EncryptionKeyResponse)
	r.WriteTo(player.Conn)

	// Start AES/CFB8 stream encryption
	player.Secret = secret
	player.OnlineMode(true)
	log.Println("Enabling encryption.")

	log.Println("HandleFC -> EncryptionKeyRequest")
}
