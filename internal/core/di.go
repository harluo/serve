package core

import (
	"github.com/harluo/di"
	"github.com/harluo/serve/internal/core/internal"
)

func init() {
	di.New().Instance().Put(
		newCommand,
		func(command *Command) internal.Put {
			return internal.Put{
				Serve: command,
			}
		},
	).Build().Apply()
}
