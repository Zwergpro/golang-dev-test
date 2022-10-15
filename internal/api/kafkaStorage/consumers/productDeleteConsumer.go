package consumers

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"google.golang.org/protobuf/proto"
	"homework-1/config"
	"homework-1/internal/cache"
	"homework-1/internal/metrics"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"time"
)

type ProductDeleteConsumer struct {
	ProductRepository repository.Product
	Metrics           *metrics.Metrics
	Cache             cache.KVCache
}

func (c *ProductDeleteConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	log.Info("starting productDeleteConsuming")
	return nil
}

func (c *ProductDeleteConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ProductDeleteConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

			in := pb.ProductDeleteRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("Failed to unmarshal message")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			err := c.ProductRepository.DeleteProduct(ctx, in.GetId())
			if err != nil {
				c.Metrics.FailedRequestCounter.Inc()
				log.WithError(err).Error("ProductRepository: DeleteProduct: internal error")
			} else {
				c.Metrics.SuccessfulRequestCounter.Inc()
				log.Infof("Product deleted: %d", in.GetId())
			}

			if err = c.Cache.Del(ctx, fmt.Sprintf("product:%d", in.GetId())); err != nil {
				log.WithError(err).Error("ProductDeleteConsumer: ConsumeClaim: del product from cache")
			}
		}
	}
}

func (c *ProductDeleteConsumer) StartConsuming(ctx context.Context) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewConsumerGroup(config.GetKafkaBrokers(), "productDeleteConsuming", saramaConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to create kafka consumer group: productDeleteConsuming")
		return
	}

	handler := otelsarama.WrapConsumerGroupHandler(c)

	for {
		if err = client.Consume(ctx, []string{"productDelete"}, handler); err != nil {
			log.WithError(err).Error("on consume productDelete")
			time.Sleep(time.Second * 3)
		}
	}
}
