//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package cli

import (
	"github.com/google/wire"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/biz"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/data"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/server"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/service"
)

func wireApp(c *conf.Config) (*server.Server, func(), error) {
	panic(wire.Build(
		biz.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
	))
}
