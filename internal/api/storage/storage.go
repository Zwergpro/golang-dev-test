package storage

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"homework-1/internal/metrics"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
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

	md, _ := metadata.FromIncomingContext(srv.Context())
	log.Infof("ProductList request metadata: %v", md)
	log.Debugf("ProductList request data: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	allProducts, err := i.deps.ProductRepository.GetAllProducts(ctx, in.GetPage(), in.GetSize())
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductRepository: GetAllProducts: internal error")
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
			log.WithError(err).Error("ProductList send")
		}
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return nil
}

func (i *implementation) ProductGet(ctx context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductGet request metadata: %v", md)
	log.Debugf("ProductGet request data: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	p, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductRepository: GetProductById: internal error")
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

func (i *implementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductCreate request metadata: %v", md)
	log.Debugf("ProductCreate request data: %v", in)

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
		log.WithError(err).Error("ProductRepository: ProductCreate: internal error")
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

func (i *implementation) ProductUpdate(ctx context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductUpdate request metadata: %v", md)
	log.Debugf("ProductUpdate request data: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	product, err := i.deps.ProductRepository.GetProductById(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductRepository: GetProductById: internal error")
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
		log.WithError(err).Error("ProductRepository: ProductUpdate: internal error")
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

func (i *implementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductDelete request metadata: %v", md)
	log.Debugf("ProductDelete request data: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err := i.deps.ProductRepository.DeleteProduct(ctx, in.GetId()); err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		if errors.Is(err, repository.ProductNotExists) {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, err.Error())
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductRepository: ProductDelete: internal error")
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pb.ProductDeleteResponse{}, nil
}
