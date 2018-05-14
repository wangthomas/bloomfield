package main

import (
	"github.com/wangthomas/bloomfield/filterManager"
	"github.com/wangthomas/bloomfield/interfaces/gRPC/gRPCHandler"
	"github.com/wangthomas/bloomfield/interfaces/gRPC/gRPCServer"

)

func main() {
	fm := filterManager.NewFilterManager()
	gh := gRPCHandler.NewgRPCHandler(fm)
	gs := gRPCServer.NewgRPCServer(gh)

	gs.Start("tcp", "8073")

}