package proxyApi

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pbApi "homework-1/pkg/api/v1"
	"log"
)

func New() pbApi.ApiServiceServer {
	return &implementation{}
}

type implementation struct {
	pbApi.UnimplementedApiServiceServer
}

func (i *implementation) ProductList(_ context.Context, in *pbApi.ProductListRequest) (*pbApi.ProductListResponse, error) {
	log.Printf("[INFO] ProductList: %v", in)
	return nil, status.Error(codes.Internal, "internal error")
}

func (i *implementation) ProductGet(_ context.Context, in *pbApi.ProductGetRequest) (*pbApi.ProductGetResponse, error) {
	log.Printf("[INFO] ProductGet: %v", in)
	return nil, status.Error(codes.Internal, "internal error")
}

func (i *implementation) ProductCreate(_ context.Context, in *pbApi.ProductCreateRequest) (*pbApi.ProductCreateResponse, error) {
	log.Printf("[INFO] ProductCreate: %v", in)
	return nil, status.Error(codes.Internal, "internal error")
}

func (i *implementation) ProductUpdate(_ context.Context, in *pbApi.ProductUpdateRequest) (*pbApi.ProductUpdateResponse, error) {
	log.Printf("[INFO] ProductUpdate: %v", in)
	return nil, status.Error(codes.Internal, "internal error")
}

func (i *implementation) ProductDelete(_ context.Context, in *pbApi.ProductDeleteRequest) (*pbApi.ProductDeleteResponse, error) {
	log.Printf("[INFO] ProductDelete: %v", in)
	return nil, status.Error(codes.Internal, "internal error")
}
