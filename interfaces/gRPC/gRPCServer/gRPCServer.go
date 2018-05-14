package gRPCServer

import (
	"net"
	
	"google.golang.org/grpc"

	pb "github.com/wangthomas/bloomfield/interfaces/gRPC/bloomfieldpb"
)

type GRPCServer struct {
	gRPCHandler pb.BloomServer
	server		*grpc.Server
}

func NewgRPCServer(gRPCHandler pb.BloomServer) *GRPCServer {
	return &GRPCServer{
		gRPCHandler:	gRPCHandler,
		server:			grpc.NewServer(),
	}
}

func (t *GRPCServer) Start(transport string, port string) error {
	listenPort := ":" + port
	listener, err := net.Listen(transport, listenPort)
	if err != nil {
		return err
	}
	pb.RegisterBloomServer(t.server, t.gRPCHandler)

	t.server.Serve(listener)

	return nil
}

