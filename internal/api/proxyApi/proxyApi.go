package proxyApi

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"io"
	"log"
	"time"
)

const maxTimeout = time.Millisecond * 30

func New(deps Deps) pbApi.ApiServiceServer {
	return &implementation{
		deps: deps,
	}
}

type implementation struct {
	pbApi.UnimplementedApiServiceServer
	deps Deps
}

type Deps struct {
	StorageClient pbStorage.StorageServiceClient
}

func (i *implementation) ProductList(_ context.Context, in *pbApi.ProductListRequest) (*pbApi.ProductListResponse, error) {
	log.Printf("[INFO] ProductList: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	pageNum := in.GetPage()
	pageSize := in.GetSize()

	productListRequest := pbStorage.ProductListRequest{Page: &pageNum, Size: &pageSize}
	productStream, err := i.deps.StorageClient.ProductList(ctx, &productListRequest)
	if err != nil {
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

	return &pbApi.ProductListResponse{
		Products: result,
	}, nil
}

func (i *implementation) ProductGet(_ context.Context, in *pbApi.ProductGetRequest) (*pbApi.ProductGetResponse, error) {
	log.Printf("[INFO] ProductGet: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	product, err := i.deps.StorageClient.ProductGet(ctx, &pbStorage.ProductGetRequest{Id: in.GetId()})
	if err != nil {
		log.Printf("[ERROR] ProductGet: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pbApi.ProductGetResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductCreate(_ context.Context, in *pbApi.ProductCreateRequest) (*pbApi.ProductCreateResponse, error) {
	log.Printf("[INFO] ProductCreate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	createProductRequest := pbStorage.ProductCreateRequest{
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	}

	product, err := i.deps.StorageClient.ProductCreate(ctx, &createProductRequest)
	if err != nil {
		log.Printf("[ERROR] ProductCreate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pbApi.ProductCreateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductUpdate(_ context.Context, in *pbApi.ProductUpdateRequest) (*pbApi.ProductUpdateResponse, error) {
	log.Printf("[INFO] ProductUpdate: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	updateProductRequest := pbStorage.ProductUpdateRequest{
		Id:       in.GetId(),
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	}

	product, err := i.deps.StorageClient.ProductUpdate(ctx, &updateProductRequest)
	if err != nil {
		log.Printf("[ERROR] ProductUpdate: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pbApi.ProductUpdateResponse{
		Id:       product.GetId(),
		Name:     product.GetName(),
		Price:    product.GetPrice(),
		Quantity: product.GetQuantity(),
	}, nil
}

func (i *implementation) ProductDelete(_ context.Context, in *pbApi.ProductDeleteRequest) (*pbApi.ProductDeleteResponse, error) {
	log.Printf("[INFO] ProductDelete: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	_, err := i.deps.StorageClient.ProductDelete(ctx, &pbStorage.ProductDeleteRequest{Id: in.GetId()})
	if err != nil {
		log.Printf("[ERROR] ProductDelete: %v\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pbApi.ProductDeleteResponse{}, nil
}
