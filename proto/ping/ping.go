package ping

import (
	"strings"

	"github.com/toqueteos/minero/constants"
	"github.com/toqueteos/minero/proto/packet"
)

// Ping returns a 0xFF packet ready to be sent (written to an io.Writer) as
// response to a 0xFE (ServerListPing) packet.
func Ping(s []string) *packet.Disconnect {
	return &packet.Disconnect{Reason: strings.Join(s, "\x00")}
}

// Prepare returns a ServerListPing-able string ready to be sent over the wire.
func Prepare(motd, in, max string) []string {
	return []string{
		"§1",
		constants.Proto, constants.Server,
		motd,
		in, max,
	}
}
