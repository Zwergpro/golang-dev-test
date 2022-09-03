package consumers

import (
	"context"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"google.golang.org/protobuf/proto"
	"homework-1/config"
	"homework-1/internal/metrics"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"time"
)

type ProductCreateConsumer struct {
	ProductRepository repository.Product
	Metrics           *metrics.Metrics
}

func (c *ProductCreateConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	log.Info("starting productCreateConsumer")
	return nil
}

func (c *ProductCreateConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ProductCreateConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			log.Info("Consume session done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Info("Data channel closed")
				return nil
			}
			c.Metrics.IncomingRequestCounter.Inc()
			session.MarkMessage(msg, "")

			in := pb.ProductCreateRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("Failed to unmarshal message")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			p := products.Product{
				Name:     in.GetName(),
				Price:    in.GetPrice(),
				Quantity: in.GetQuantity(),
			}

			product, err := c.ProductRepository.CreateProduct(ctx, p)
			if err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("ProductRepository: ProductCreate: internal error")
			} else {
				c.Metrics.SuccessfulRequestCounter.Inc()
				log.Infof("Product created: %v", product)
			}
		}
	}
}

func (c *ProductCreateConsumer) StartConsuming(ctx context.Context) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewConsumerGroup(config.GetKafkaBrokers(), "productCreateConsuming", saramaConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to create kafka consumer group: productCreateConsuming")
	}

	handler := otelsarama.WrapConsumerGroupHandler(c)

	for {
		if err := client.Consume(ctx, []string{"productCreate"}, handler); err != nil {
			log.WithError(err).Error("on consume productCreate")
			time.Sleep(time.Second * 3)
		}
	}
}
