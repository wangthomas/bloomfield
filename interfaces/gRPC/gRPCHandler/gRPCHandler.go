package gRPCHandler

import (
	"context"

	fm "github.com/wangthomas/bloomfield/filterManager"
	pb "github.com/wangthomas/bloomfield/interfaces/gRPC/bloomfieldpb"
	
)



type gRPCHandler struct {
	filterManager *fm.FilterManager
}

func NewgRPCHandler(filterManager *fm.FilterManager) (pb.BloomServer) {
	return &gRPCHandler{
		filterManager: filterManager,
	}
}

func (t *gRPCHandler) CreateFilter(ctx context.Context, request *pb.FilterRequest) (*pb.Response, error) {

	t.filterManager.Create(request.Name)

	return &pb.Response{
		Status: pb.Status_SUCCESS,
	}, nil
}


func (t *gRPCHandler) DropFilter(ctx context.Context, request *pb.FilterRequest) (*pb.Response, error) {
	t.filterManager.Drop(request.Name)

	return &pb.Response{
		Status: pb.Status_SUCCESS,
	}, nil
}

func (t *gRPCHandler) Add(ctx context.Context, request *pb.KeyRequest) (*pb.HasResponse, error) {
	
	return &pb.HasResponse{
		Status: pb.Status_SUCCESS,
		Has:	t.filterManager.Add(request.FilterName, request.Hashes),
	}, nil
}

func (t *gRPCHandler) Has(ctx context.Context, request *pb.KeyRequest) (*pb.HasResponse, error) {

	return &pb.HasResponse{
		Status: pb.Status_SUCCESS,
		Has:	t.filterManager.Has(request.FilterName, request.Hashes),
	}, nil

}



