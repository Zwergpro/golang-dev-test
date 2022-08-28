package consumers

import (
	"context"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"time"
)

type ProductDeleteConsumer struct {
	ProductRepository repository.Product
}

func (c *ProductDeleteConsumer) Setup(_ sarama.ConsumerGroupSession) error {
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

			in := pb.ProductDeleteRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				log.WithError(err).Error("Failed to unmarshal message")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			err := c.ProductRepository.DeleteProduct(ctx, in.GetId())
			if err != nil {
				log.WithError(err).Error("ProductRepository: DeleteProduct: internal error")
			} else {
				log.Infof("Product deleted: %d", in.GetId())
			}
			session.MarkMessage(msg, "")
		}
	}
}
