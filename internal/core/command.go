package core

import (
	"context"
	"sync"
	"time"

	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/harluo/boot"
	"github.com/harluo/serve/internal/kernel"
)

type Command struct {
	servers []kernel.Server

	logger  log.Logger
	exiting bool
}

func newCommand(logger log.Logger) *Command {
	return &Command{
		servers: make([]kernel.Server, 0),

		logger: logger,
	}
}

func (c *Command) Add(required kernel.Server, optionals ...kernel.Server) (command *Command) {
	c.servers = append(c.servers, required)
	c.servers = append(c.servers, optionals...)
	command = c

	return
}

func (c *Command) Name() string {
	return "serve"
}

func (c *Command) Aliases() []string {
	return []string{
		"s",
		"srv",
	}
}

func (c *Command) Usage() string {
	return "启动服务"
}

func (c *Command) Hidden() bool {
	return false
}

func (c *Command) Run(ctx context.Context) (err error) {
	count := len(c.servers)
	if 0 != count {
		c.logger.Debug("启动所有服务开始", field.New("count", count))
		err = c.start(ctx, count)
	}

	return
}

func (c *Command) Stop(ctx context.Context) (err error) {
	c.exiting = true
	wg := new(sync.WaitGroup)
	wg.Add(len(c.servers))
	for _, serve := range c.servers {
		go c.stopServer(ctx, serve, wg, &err)
	}
	wg.Wait()

	return
}

func (c *Command) Before(ctx context.Context) (err error) {
	/*for _, server := range c.servers {
		err = server.Before(ctx)
		if nil != err {
			break
		}
	}*/

	return
}

func (c *Command) After(ctx context.Context) (err error) {
	/*for _, server := range c.servers {
		err = server.After(ctx)
		if nil != err {
			break
		}
	}*/

	return
}

func (c *Command) Arguments() boot.Arguments {
	return boot.Arguments{}
}

func (c *Command) Subcommands() boot.Commands {
	return boot.Commands{}
}

func (c *Command) Description() string {
	return ""
}

func (c *Command) Category() string {
	return ""
}

func (c *Command) start(ctx context.Context, count int) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(count)
	for _, serve := range c.servers {
		cloned := serve
		go c.startServer(ctx, cloned, wg, &err)
	}
	c.logger.Debug("启动所有服务成功", field.New("count", count))
	wg.Wait()

	return
}

func (c *Command) startServer(ctx context.Context, server kernel.Server, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	c.logger.Info("启动服务成功", field.New[string]("name", server.Name()))
	// 记录时间，如果发生错误的时间小于500毫秒，就是执行错误，应该立即退出；如果大于，则只记录日志
	now := time.Now()
	if se := server.Start(ctx); nil != se && !c.exiting {
		errTime := time.Now()
		if errTime.Sub(now) > 500*time.Millisecond {
			c.logger.Info("服务执行错误", field.New[string]("name", server.Name()), field.Error(se))
		} else {
			*err = se
		}
	}
}

func (c *Command) stopServer(ctx context.Context, server kernel.Server, wg *sync.WaitGroup, err *error) {
	defer wg.Done()
	if se := server.Stop(ctx); nil != se {
		*err = se
		c.logger.Info("停止服务出错", field.New("name", server.Name()), field.Error(se))
	} else {
		c.logger.Info("停止服务成功", field.New("name", server.Name()))
	}
}
