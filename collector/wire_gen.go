// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package collector

import (
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/geo"
)

// Injectors from wire.go:

func NewCollector() (*Collector, error) {
	configurationConfiguration := configuration.GetConfig()
	connection := db.Connect(configurationConfiguration)
	geoConnection := geo.Connect(configurationConfiguration)
	inputParams := initInputParams()
	collector := &Collector{
		db:          connection,
		geo:         geoConnection,
		config:      configurationConfiguration,
		inputParams: inputParams,
	}
	return collector, nil
}