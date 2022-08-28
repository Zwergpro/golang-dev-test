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

type ProductUpdateConsumer struct {
	ProductRepository repository.Product
}

func (c *ProductUpdateConsumer) Setup(_ sarama.ConsumerGroupSession) error {
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

			in := pb.ProductUpdateRequest{}
			if err := proto.Unmarshal(msg.Value, &in); err != nil {
				log.WithError(err).Error("Failed to unmarshal message")
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			product, err := c.ProductRepository.GetProductById(ctx, in.GetId())
			if err != nil {
				log.WithError(err).Error("ProductRepository: GetProductById: internal error")
				continue
			}

			product.Name = in.GetName()
			product.Price = in.GetPrice()
			product.Quantity = in.GetQuantity()

			if product, err = c.ProductRepository.UpdateProduct(ctx, *product); err != nil {
				log.WithError(err).Error("ProductRepository: ProductUpdate: internal error")
			} else {
				log.Infof("Product updated: %v", product)
			}

			session.MarkMessage(msg, "")
		}
	}
}
