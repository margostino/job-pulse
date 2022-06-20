package main

import (
	"github.com/margostino/job-pulse/collector"
	"github.com/margostino/job-pulse/utils"
)

//const (
//	SearchPosition = "software engineer"
//	SearchLocation = "stockholm"
//)

var searchPosition = ""
var searchLocation = ""

func main() {
	collector, err := collector.NewCollector()
	utils.Check(err)
	collector.Start()
	collector.Close()
}
