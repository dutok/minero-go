// Package minecraft implements server list ping of Minecraft Protocol.
//
// Many, many thanks to #mcdevs from Freenode and it's great documentation:
// http://wiki.vg/Server_List_Ping
package ping

import (
	"strings"
	"unicode/utf16"

	"github.com/toqueteos/minero/proto/packet"
)

// Ping returns a 0xFF packet containing the response to a 0xFE (ServerListPing)
// packet. For more info check package docs.
func Ping(s []string) *packet.Disconnect {
	return &packet.Disconnect{Reason: strings.Join(s, "\x00")}
}
