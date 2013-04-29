package vanilla

import (
	"github.com/toqueteos/minero/cmd"
	vcmd "github.com/toqueteos/minero/vanilla/cmd"
)

var CmdList = map[string]cmd.Cmder{
	"gamemode": new(vcmd.Gamemode),
}
