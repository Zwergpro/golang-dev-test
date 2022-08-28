package consumers

import (
	"context"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"google.golang.org/protobuf/proto"
	"homework-1/config"
	"homework-1/internal/metrics"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"time"
)

type ProductUpdateConsumer struct {
	ProductRepository repository.Product
	Metrics           *metrics.Metrics
}

func (c *ProductUpdateConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	log.Info("starting productUpdateConsumer")
	return nil
}

func (c *ProductUpdateConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ProductUpdateConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

			in := pb.ProductUpdateRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("Failed to unmarshal message")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			product, err := c.ProductRepository.GetProductById(ctx, in.GetId())
			if err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("ProductRepository: GetProductById: internal error")
				continue
			}

			product.Name = in.GetName()
			product.Price = in.GetPrice()
			product.Quantity = in.GetQuantity()

			if product, err = c.ProductRepository.UpdateProduct(ctx, *product); err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("ProductRepository: ProductUpdate: internal error")
			} else {
				c.Metrics.SuccessfulRequestCounter.Inc()
				log.Infof("Product updated: %v", product)
			}
		}
	}
}

func (c *ProductUpdateConsumer) StartConsuming(ctx context.Context) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewConsumerGroup(config.GetKafkaBrokers(), "productUpdateConsumer", saramaConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to create kafka consumer group: productUpdateConsumer")
		return
	}

	handler := otelsarama.WrapConsumerGroupHandler(c)

	for {
		if err := client.Consume(ctx, []string{"productUpdate"}, handler); err != nil {
			log.WithError(err).Error("on consume productUpdate")
			time.Sleep(time.Second * 3)
		}
	}
}
