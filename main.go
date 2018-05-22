package main

import (
	"fmt"
	"flag"

	"github.com/wangthomas/bloomfield/filterManager"
	"github.com/wangthomas/bloomfield/config"
	"github.com/wangthomas/bloomfield/interfaces/gRPC/gRPCHandler"
	"github.com/wangthomas/bloomfield/interfaces/gRPC/gRPCServer"

)

func main() {

	var configFile string

	flag.StringVar(&configFile, "config_file", "", "specify config file")

	flag.Parse()

	if configFile != "" {
		config.LoadConfig(configFile)
	} else {
		fmt.Println("Using default config")
		config.LoadDefault()
	}

	filtermanager := filterManager.NewFilterManager()
	ghander := gRPCHandler.NewgRPCHandler(filtermanager)
	gserver := gRPCServer.NewgRPCServer(ghander)

	gserver.Start("tcp", config.Config.Port)

}