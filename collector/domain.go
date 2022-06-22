package collector

import (
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/geo"
	"github.com/margostino/job-pulse/scrapper"
)

type InputParams struct {
	SearchPosition string
	SearchLocation string
}

type App struct {
	inputParams *InputParams
	db          *db.Connection
	geo         *geo.Connection
	scrapper    *scrapper.Scrapper
	config      *configuration.Configuration
}
