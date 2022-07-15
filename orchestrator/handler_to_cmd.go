package orchestrator

import (
	"kubinka/cmd_deploy"
	"kubinka/cmd_return"
	"kubinka/command"
)

var handlerToCmd = map[string]func() command.Command{
	"deploy": cmd_deploy.Init,
	"return": cmd_return.Init,
}
