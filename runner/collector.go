package main

import (
	"github.com/margostino/job-pulse/collector"
	"github.com/margostino/job-pulse/utils"
)

func main() {
	app, err := collector.NewApp()
	utils.Check(err)
	app.ValidateInput()
	app.Start()
	app.Close()
}
