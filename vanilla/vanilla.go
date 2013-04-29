package vanilla

import (
	"github.com/toqueteos/minero/cmd"
)

var CmdList = map[string]cmd.Cmder{
	"gamemode": new(Gamemode),
}
