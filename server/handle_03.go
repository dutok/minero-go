package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/toqueteos/minero/cmd"
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle03 handles incoming requests of packet 0x03: ChatMessage
func Handle03(server *Server, sender *player.Player) {
	pkt := new(packet.ChatMessage)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ChatMessage: %+v", pkt)

	// Messages prefixed with / are treated like commands
	if strings.HasPrefix(pkt.Message, "/") {
		var parts = strings.Fields(pkt.Message[1:])
		// Empty commands are noops.
		if len(parts) == 0 {
			return
		}

		var cmdName, cmdArgs = parts[0], parts[1:]
		var cmd cmd.Cmder
		var ok bool

		// Command not found
		if cmd, ok = server.Cmds[cmdName]; !ok {
			msg := fmt.Sprintf("Unknown command %q.", cmdName)
			log.Println(msg)
			sender.SendMessage(msg)
			return
		}

		ok = cmd.Do(sender, cmdArgs)
		if !ok {
			msg := "An error ocurred executing command %q."
			log.Println(msg)
			sender.SendMessage(msg)
			return
		}
	}

	// Send message to all other players
	msg := fmt.Sprintf("<%s> %s", sender.Name, pkt.Message)
	sender.BroadcastMessage(server.Players.Copy(), msg)
}

func contains(cmdName string, list map[string]cmd.Cmder) (ok bool) {
	_, ok = list[cmdName]
	return
}
