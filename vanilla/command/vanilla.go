package vanilla

import (
	"github.com/toqueteos/minero/command"
)

var CmdList = map[string]command.Cmder{
	"gamemode": Gamemode{},
}
