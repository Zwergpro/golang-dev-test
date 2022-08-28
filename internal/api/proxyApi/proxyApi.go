//go:generate mockgen -source ./proxyApi.go -destination=./mock/storage.go -package=mock_storage

package proxyApi

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/metrics"
	"homework-1/internal/models/products"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"io"
	"log"
	"strings"
	"time"
)

const maxTimeout = time.Millisecond * 30

type StorageServiceClient interface {
	pbStorage.StorageServiceClient
}

func New(deps Deps) *implementation {
	return &implementation{
		deps: deps,
	}
}

type implementation struct {
	pbApi.UnimplementedApiServiceServer
	deps Deps
}

type Deps struct {
	StorageClient StorageServiceClient
	Metrics       *metrics.Metrics
}

func (i *implementation) ProductList(_ context.Context, in *pbApi.ProductListRequest) (*pbApi.ProductListResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductList: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	pageNum := in.GetPage()
	pageSize := in.GetSize()

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	request := pbStorage.ProductListRequest{Page: &pageNum, Size: &pageSize}
	productStream, err := i.deps.StorageClient.ProductList(ctx, &request)
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductList: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var result []*pbApi.ProductListResponse_Product
	for {
		product, err := productStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			i.deps.Metrics.FailedRequestCounter.Inc()
			log.Printf("[ERROR] ProductList receive: %v\n", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
		result = append(result, &pbApi.ProductListResponse_Product{
			Id:       product.GetId(),
			Name:     product.GetName(),
			Price:    product.GetPrice(),
			Quantity: product.GetQuantity(),
		})
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductListResponse{
		Products: result,
	}, nil
}

func (i *implementation) ProductGet(_ context.Context, in *pbApi.ProductGetRequest) (*pbApi.ProductGetResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductGet: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	product, err := i.deps.StorageClient.ProductGet(ctx, &pbStorage.ProductGetRequest{Id: in.GetId()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, "product not found")
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductGet: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductGetResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductCreate(_ context.Context, in *pbApi.ProductCreateRequest) (*pbApi.ProductCreateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductCreate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if errs := products.ValidateProductFields(in.GetName(), in.GetPrice(), in.GetQuantity()); len(errs) > 0 {
		errStrings := make([]string, 0, len(errs))
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
		return nil, status.Error(codes.InvalidArgument, strings.Join(errStrings, "; "))
	}

	request := pbStorage.ProductCreateRequest{
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	}

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	product, err := i.deps.StorageClient.ProductCreate(ctx, &request)
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductCreate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductCreateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pbApi.ProductUpdateRequest) (*pbApi.ProductUpdateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductUpdate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if errs := products.ValidateProductFields(in.GetName(), in.GetPrice(), in.GetQuantity()); len(errs) > 0 {
		errStrings := make([]string, 0, len(errs))
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
		return nil, status.Error(codes.InvalidArgument, strings.Join(errStrings, "; "))
	}

	request := pbStorage.ProductUpdateRequest{
		Id:       in.GetId(),
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	}

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	product, err := i.deps.StorageClient.ProductUpdate(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, "product not found")
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductUpdateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductDelete(_ context.Context, in *pbApi.ProductDeleteRequest) (*pbApi.ProductDeleteResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()
	log.Printf("[INFO] ProductDelete: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	_, err := i.deps.StorageClient.ProductDelete(ctx, &pbStorage.ProductDeleteRequest{Id: in.GetId()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			i.deps.Metrics.UnsuccessfulRequestCounter.Inc()
			return nil, status.Error(codes.NotFound, "product not found")
		}
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.Printf("[ERROR] ProductDelete: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductDeleteResponse{}, nil
}
