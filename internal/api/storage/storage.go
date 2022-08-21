package storage

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"log"
	"time"
)

const maxTimeout = time.Millisecond * 27

func New(deps Deps) *implementation {
	return &implementation{
		deps: deps,
	}
}

type implementation struct {
	pb.UnimplementedStorageServiceServer
	deps Deps
}

type Deps struct {
	ProductRepository repository.Product
}

func (i *implementation) ProductList(in *pb.ProductListRequest, srv pb.StorageService_ProductListServer) error {
	log.Printf("[INFO] ProductList: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	allProducts, err := i.deps.ProductRepository.GetAllProducts(ctx, in.GetPage(), in.GetSize())
	if err != nil {
		log.Printf("[ERROR] ProductList: %v\n", err)
		return status.Error(codes.Internal, "internal error")
	}

	for _, product := range allProducts {
		productResponse := pb.ProductListResponse{
			Id:       product.GetId(),
			Name:     product.GetName(),
			Price:    product.GetPrice(),
			Quantity: product.GetQuantity(),
		}
		if err = srv.Send(&productResponse); err != nil {
			log.Printf("[ERROR] ProductList send: %v\n", product)
		}
	}

	return nil
}

func (i *implementation) ProductGet(_ context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	log.Printf("[INFO] ProductGet: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[ERROR] ProductGet: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.ProductGetResponse{
		Id:       p.GetId(),
		Name:     p.GetName(),
		Price:    p.GetPrice(),
		Quantity: p.GetQuantity(),
	}, nil
}

func (i *implementation) ProductCreate(_ context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	log.Printf("[INFO] ProductCreate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := products.BuildProduct(in.GetName(), in.GetPrice(), in.GetQuantity())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	product, err := i.deps.ProductRepository.CreateProduct(ctx, *p)
	if err != nil {
		if errors.As(err, &repository.ProductAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		log.Printf("[ERROR] ProductCreate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.ProductCreateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	log.Printf("[INFO] ProductUpdate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	product, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.As(err, &repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	if err = product.SetName(in.GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = product.SetPrice(in.GetPrice()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = product.SetQuantity(in.GetQuantity()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if product, err = i.deps.ProductRepository.UpdateProduct(ctx, *product); err != nil {
		if errors.As(err, &repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.ProductUpdateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductDelete(_ context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	log.Printf("[INFO] ProductDelete: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := i.deps.ProductRepository.DeleteProduct(ctx, in.GetId()); err != nil {
		if errors.As(err, &repository.ProductNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[ERROR] ProductDelete: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.ProductDeleteResponse{}, nil
}
