//go:build wireinject
// +build wireinject

package collector

import (
	"github.com/google/wire"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/geo"
	"github.com/margostino/job-pulse/scrapper"
)

func NewApp() (*App, error) {
	wire.Build(
		db.Connect,
		geo.Connect,
		scrapper.New,
		configuration.GetConfig,
		NewInputParams,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
