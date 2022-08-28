package storage

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/metrics"
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
	Metrics           *metrics.Metrics
}

func (i *implementation) ProductList(in *pb.ProductListRequest, srv pb.StorageService_ProductListServer) error {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductList: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	allProducts, err := i.deps.ProductRepository.GetAllProducts(ctx, in.GetPage(), in.GetSize())
	if err != nil {
		log.Printf("[ERROR] ProductList: %v\n", err)
		i.deps.Metrics.FailedRequestCounter.Inc()
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

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return nil
}

func (i *implementation) ProductGet(_ context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductGet: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductGet: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pb.ProductGetResponse{
		Id:       p.GetId(),
		Name:     p.GetName(),
		Price:    p.GetPrice(),
		Quantity: p.GetQuantity(),
	}, nil
}

func (i *implementation) ProductCreate(_ context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductCreate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p := products.Product{
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	}

	product, err := i.deps.ProductRepository.CreateProduct(ctx, p)
	if err != nil {
		if errors.Is(err, repository.ProductAlreadyExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductCreate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pb.ProductCreateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductUpdate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	product, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	product.Name = in.GetName()
	product.Price = in.GetPrice()
	product.Quantity = in.GetQuantity()

	if product, err = i.deps.ProductRepository.UpdateProduct(ctx, *product); err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pb.ProductUpdateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductDelete(_ context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductDelete: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := i.deps.ProductRepository.DeleteProduct(ctx, in.GetId()); err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductDelete: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pb.ProductDeleteResponse{}, nil
}
