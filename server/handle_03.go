package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/toqueteos/minero/command"
	"github.com/toqueteos/minero/command/vanilla"
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
		var (
			parts            = strings.Fields(pkt.Message[1:])
			cmdName, cmdArgs = parts[0], parts[1:]
			cmd              command.Cmder
			ok               bool
		)

		switch {
		case contains(cmdName, server.cmdList):
			cmd = vanilla.CmdList[cmdName]
		case contains(cmdName, vanilla.CmdList):
			cmd = vanilla.CmdList[cmdName]
		default:
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

	// All other messages are sent
	msg := fmt.Sprintf("<%s> %s", sender.Name, pkt.Message)
	server.BroadcastMessage(msg)
}

func contains(cmdName string, list map[string]command.Cmder) (ok bool) {
	_, ok = list[cmdName]
	return
}
