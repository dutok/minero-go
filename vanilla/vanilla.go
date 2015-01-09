package vanilla

import (
	"github.com/minero/minero/cmd"
	vcmd "github.com/minero/minero/vanilla/cmd"
)

var CmdList = map[string]cmd.Cmder{
	"gamemode": new(vcmd.Gamemode),
}
