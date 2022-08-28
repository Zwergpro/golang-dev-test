package consumers

import (
	"context"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	pb "homework-1/pkg/api/storage/v1"
	"time"
)

type ProductCreateConsumer struct {
	ProductRepository repository.Product
}

func (c *ProductCreateConsumer) Setup(_ sarama.ConsumerGroupSession) error {
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

			in := pb.ProductCreateRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				log.WithError(err).Error("Failed to unmarshal message")
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
				log.WithError(err).Error("ProductRepository: ProductCreate: internal error")
			} else {
				log.Infof("Product created: %v", product)
			}
			session.MarkMessage(msg, "")
		}
	}
}
