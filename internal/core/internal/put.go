package internal

import (
	"github.com/harluo/boot"
	"github.com/harluo/di"
)

type Put struct {
	di.Put

	Serve boot.Command `group:"commands"`
}
