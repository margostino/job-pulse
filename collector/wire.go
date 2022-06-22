//go:build wireinject
// +build wireinject

package collector

import (
	"github.com/google/wire"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/geo"
)

func NewCollector() (*Collector, error) {
	wire.Build(
		db.Connect,
		geo.Connect,
		configuration.GetConfig,
		newInputParams,
		wire.Struct(new(Collector), "*"),
	)
	return &Collector{}, nil
}
