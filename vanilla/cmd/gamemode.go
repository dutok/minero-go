package cmd

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

type Gamemode struct{}

func (g *Gamemode) Tab(args []string) []string { return nil }

func (g *Gamemode) Do(from *player.Player, args []string) bool {
	r := packet.GameStateChange{Reason: 3}

	switch arg := args[0]; arg {
	case "survival":
		r.GameMode = 0
	case "creative":
		r.GameMode = 1
	case "adventure":
		r.GameMode = 2
	default:
		// First char of first arg
		switch arg[0] {
		case '0', 's':
			r.GameMode = 0
		case '1', 'c':
			r.GameMode = 1
		case '2', 'a':
			r.GameMode = 2
		}
	}

	_, err := r.WriteTo(from.Conn)

	return err != nil
}
