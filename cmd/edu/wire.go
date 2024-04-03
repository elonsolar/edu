//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"edu/internal/conf"
	"edu/internal/domain"
	"edu/internal/infra/repo"
	"edu/internal/server"
	"edu/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, repo.ProviderSet, domain.ProviderSet, service.ProviderSet, newApp))
}
