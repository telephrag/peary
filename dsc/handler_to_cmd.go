package dsc

import (
	"context"
	"kubinka/deploy"
)

var HandlerToCmd = map[string]func(context.Context) interface{}{
	"deploy": deploy.Init,
	"return": nil,
}
