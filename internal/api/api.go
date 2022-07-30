package api

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/models"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/v1"
	"time"
)

const maxTimeout = time.Millisecond * 27

func New(deps Deps) pb.AdminServiceServer {
	return &implementation{
		deps: deps,
	}
}

type implementation struct {
	pb.UnimplementedAdminServiceServer
	deps Deps
}

type Deps struct {
	ProductRepository repository.Product
}

func (i *implementation) ProductList(_ context.Context, _ *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	products, err := i.deps.ProductRepository.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.As(err, &repository.ProductNotExists) {
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
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := models.BuildProduct(in.GetName(), in.GetPrice(), in.GetQuantity())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	product, err := i.deps.ProductRepository.CreateProduct(ctx, *p)
	if err != nil {
		if errors.As(err, &repository.ProductAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductCreateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.As(err, &repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = p.SetName(in.GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = p.SetPrice(in.GetPrice()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = p.SetQuantity(in.GetQuantity()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if p, err = i.deps.ProductRepository.UpdateProduct(ctx, *p); err != nil {
		if errors.As(err, &repository.ProductNotExists) {
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
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := i.deps.ProductRepository.DeleteProduct(ctx, in.GetId()); err != nil {
		if errors.As(err, &repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ProductDeleteResponse{}, nil
}
