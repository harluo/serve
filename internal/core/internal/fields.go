package internal

import (
	"github.com/goexl/gox"
)

type Fields interface {
	Fields() []gox.Field[any]
}
