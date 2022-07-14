package dsc

import (
	"kubinka/cmd_return"
	"kubinka/command"
	cmd_deploy "kubinka/deploy"
)

var HandlerToCmd = map[string]func() command.Command{
	"deploy": cmd_deploy.Init,
	"return": cmd_return.Init,
}
