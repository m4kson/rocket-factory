package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/m4kson/rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/m4kson/rocket-factory/assembly/internal/converter/kafka"
	"github.com/m4kson/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/m4kson/rocket-factory/assembly/internal/service"
	orderConsumer "github.com/m4kson/rocket-factory/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/m4kson/rocket-factory/assembly/internal/service/producer/order_producer"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/m4kson/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/m4kson/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/m4kson/rocket-factory/platform/pkg/kafka/producer"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	kafkaMiddleware "github.com/m4kson/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	shipAssembledProducerService service.ShipAssembledProducerService
	orderPaidConsumerService     service.ConsumerService

	consumerGroup         sarama.ConsumerGroup
	orderRecordedConsumer wrappedKafka.Consumer
	orderRecordedDecoder  kafkaConverter.OrderPaidDecoder

	syncProducer          sarama.SyncProducer
	orderRecordedProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ShipAssembledProducerService(ctx context.Context) service.ShipAssembledProducerService {
	if d.shipAssembledProducerService == nil {
		d.shipAssembledProducerService = orderProducer.NewService(d.OrderRecordedProducer(ctx))
	}

	return d.shipAssembledProducerService
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderConsumer.NewService(d.OrderRecordedConsumer(ctx), d.OrderRecordedDecoder(), d.ShipAssembledProducerService(ctx))
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)

		if err != nil {
			fmt.Sprintf("failed to create consumer group: %s\n", err.Error())
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderRecordedConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderRecordedConsumer == nil {
		d.orderRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.FromContext(ctx),
			kafkaMiddleware.Logging(logger.FromContext(ctx)),
		)
	}

	return d.orderRecordedConsumer
}

func (d *diContainer) OrderRecordedDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderRecordedDecoder == nil {
		d.orderRecordedDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderRecordedDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderRecordedProducer(ctx context.Context) wrappedKafka.Producer {
	if d.orderRecordedProducer == nil {
		d.orderRecordedProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.FromContext(ctx),
		)
	}

	return d.orderRecordedProducer
}
