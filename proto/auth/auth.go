package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
)

// GenerateKeyPair generates a 1024 bit RSA keypair.
func GenerateKeyPair() *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil
	}
	return priv
}

// KeyExchange marshals a RSA Public Key
func KeyExchange(pub *rsa.PublicKey) []byte {
	asn1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}
	return asn1
}

// EncryptionBytes returns 4 random bytes or nil. Useful for Protocol Encryption
// (packets 0xFC & 0xFD).
func EncryptionBytes() []byte {
	var buf bytes.Buffer
	n, err := io.CopyN(&buf, rand.Reader, 4)
	if n != 4 || err != nil {
		return nil
	}
	return buf.Bytes()
}
