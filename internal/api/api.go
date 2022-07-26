package api

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/storage"
	pb "homework-1/pkg/api/v1"
)

func New() pb.AdminServiceServer {
	return &implementation{}
}

type implementation struct {
	pb.UnimplementedAdminServiceServer
}

func (i *implementation) ProductList(_ context.Context, _ *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	products := storage.List()

	result := make([]*pb.ProductListResponse_Product, 0, len(products))
	for _, product := range products {
		result = append(result, &pb.ProductListResponse_Product{
			Id:       product.GetId(),
			Name:     product.GetName(),
			Price:    product.GetPrice(),
			Quantity: product.GetQuantity(),
		})
	}

	return &pb.ProductListResponse{
		Products: result,
	}, nil
}

func (i *implementation) ProductGet(_ context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	p, err := storage.Get(in.GetId())
	if err != nil {
		if errors.As(err, &storage.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ProductGetResponse{
		Id:       p.GetId(),
		Name:     p.GetName(),
		Price:    p.GetPrice(),
		Quantity: p.GetQuantity(),
	}, nil
}

func (i *implementation) ProductCreate(_ context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	p, err := storage.NewProduct(in.GetName(), in.GetPrice(), in.GetQuantity())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = storage.Add(p); err != nil {
		if errors.As(err, &storage.ProductAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductCreateResponse{
		Id:       p.GetId(),
		Name:     p.GetName(),
		Price:    p.GetPrice(),
		Quantity: p.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	oldProduct, err := storage.Get(in.GetId())
	if err != nil {
		if errors.As(err, &storage.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	p := oldProduct.Copy()

	if err = p.SetName(in.GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = p.SetPrice(in.GetPrice()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = p.SetQuantity(in.GetQuantity()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = storage.Update(p); err != nil {
		if errors.As(err, &storage.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductUpdateResponse{
		Id:       p.GetId(),
		Name:     p.GetName(),
		Price:    p.GetPrice(),
		Quantity: p.GetQuantity(),
	}, nil
}

func (i *implementation) ProductDelete(_ context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	if err := storage.Delete(in.GetId()); err != nil {
		if errors.As(err, &storage.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ProductDeleteResponse{}, nil
}
