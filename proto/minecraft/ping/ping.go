package ping

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

// Ping prepares an 0xFF packet
func Ping(ping []string) []byte {
	// Write Disconnect packet Id (byte)
	buf := bytes.NewBuffer([]byte{0xff})

	// Write length of string (short)
	payload := prepare(ping)
	length := uint16(len(payload) / 2)      // 2 bytes
	buf.WriteByte(byte(length >> 8 & 0xFF)) // left byte
	buf.WriteByte(byte(length & 0xFF))      // right byte
	buf.Write(payload)

	return buf.Bytes()
}

func prepare(ping []string) []byte {
	var payload bytes.Buffer

	s := strings.Join(ping, "\x00")
	ucs2 := utf16.Encode([]rune(s))
	err := binary.Write(&payload, binary.BigEndian, ucs2)
	if err != nil {
		return nil
	}

	return payload.Bytes()
}
