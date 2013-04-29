package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleFC handles incoming requests of packet 0xFC: EncryptionKeyResponse
func HandleFC(server *Server, sender *player.Player) {
	pkt := new(packet.EncryptionKeyResponse)
	pkt.ReadFrom(sender.Conn)

	// Decrypt shared secret and token with server's private key.
	var secret, token []byte
	// var err error
	secret, _ = rsa.DecryptPKCS1v15(rand.Reader, server.PrivateKey(), pkt.Secret)
	token, _ = rsa.DecryptPKCS1v15(rand.Reader, server.PrivateKey(), pkt.Token)

	// Ensure token matches
	if !bytes.Equal(token, sender.Token) {
		log.Println("Tokens don't match.")
		r := &packet.Disconnect{Reason: ReasonPiratedGame}
		r.WriteTo(sender.Conn)
		return
	}

	// Ensure player is legit
	if !server.CheckUser(sender.Name, secret) {
		log.Println("Failed to verify username!")
		r := packet.Disconnect{"Failed to verify username!"}
		r.WriteTo(sender.Conn)
		return
	}

	// Send empty EncryptionKeyResponse
	r := new(packet.EncryptionKeyResponse)
	r.WriteTo(sender.Conn)

	// Start AES/CFB8 stream encryption
	sender.OnlineMode(true, secret)
	log.Println("Enabling encryption.")
}
