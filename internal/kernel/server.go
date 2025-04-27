package kernel

import (
	"context"
)

type Server interface {
	// Start 运行服务
	Start(ctx context.Context) error

	// Stop 停止
	Stop(ctx context.Context) error

	// Name 服务名称
	Name() string
}
