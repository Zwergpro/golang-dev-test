package kafkaProxyApi

import (
	"context"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"homework-1/internal/metrics"
	"homework-1/internal/models/products"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"io"
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
	Producer      sarama.SyncProducer
	Metrics       *metrics.Metrics
}

func (i *implementation) ProductList(ctx context.Context, in *pbApi.ProductListRequest) (*pbApi.ProductListResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductList request metadata: %v", md)
	log.Debugf("ProductList request data: %v", in)

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	pageNum := in.GetPage()
	pageSize := in.GetSize()

	i.deps.Metrics.OutgoingRequestCounter.Inc()
	request := pbStorage.ProductListRequest{Page: &pageNum, Size: &pageSize}
	productStream, err := i.deps.StorageClient.ProductList(ctx, &request)
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("StorageClient: ProductList: internal error")
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
			log.WithError(err).Error("StorageClient: ProductList: receive internal error")
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

func (i *implementation) ProductGet(ctx context.Context, in *pbApi.ProductGetRequest) (*pbApi.ProductGetResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductGet request metadata: %v", md)
	log.Debugf("ProductGet request data: %v", in)

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
		log.WithError(err).Error("StorageClient: ProductGet: internal error")
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

func (i *implementation) ProductCreate(ctx context.Context, in *pbApi.ProductCreateRequest) (*pbApi.ProductCreateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductCreate request metadata: %v", md)
	log.Debugf("ProductCreate request data: %v", in)

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

	requestData, err := proto.Marshal(&pbStorage.ProductCreateRequest{
		Name:     in.GetName(),
		Price:    in.GetPrice(),
		Quantity: in.GetQuantity(),
	})
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductCreate: proto.Marshal: internal error")
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.OutgoingRequestCounter.Inc()

	_, _, err = i.deps.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "productCreate",
		Value: sarama.ByteEncoder(requestData),
	})
	if err != nil {
		i.deps.Metrics.FailedRequestCounter.Inc()
		log.WithError(err).Error("ProductCreate: Producer: SendMessage: internal error")
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductCreateResponse{
		Id:       1,
		Name:     "test",
		Price:    1,
		Quantity: 1,
	}, nil
}

func (i *implementation) ProductUpdate(ctx context.Context, in *pbApi.ProductUpdateRequest) (*pbApi.ProductUpdateResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductUpdate request metadata: %v", md)
	log.Debugf("ProductUpdate request data: %v", in)

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
		log.WithError(err).Error("StorageClient: ProductUpdate: internal error")
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

func (i *implementation) ProductDelete(ctx context.Context, in *pbApi.ProductDeleteRequest) (*pbApi.ProductDeleteResponse, error) {
	i.deps.Metrics.IncomingRequestCounter.Inc()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Infof("ProductDelete request metadata: %v", md)
	log.Debugf("ProductDelete request data: %v", in)

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
		log.WithError(err).Error("StorageClient: ProductDelete: internal error")
		return nil, status.Error(codes.Internal, "internal error")
	}

	i.deps.Metrics.SuccessfulRequestCounter.Inc()
	return &pbApi.ProductDeleteResponse{}, nil
}
