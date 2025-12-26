package bootstrap

import (
	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/consumer"
	"github.com/qnhqn1/file-validator/internal/metrics"
	"github.com/qnhqn1/file-validator/internal/producer"
	"github.com/qnhqn1/file-validator/internal/services/validator"
)


func InitConsumers(cfg *config.Config, service validator.Service, collector *metrics.Collector, producers *producer.Manager) *consumer.Manager {
	return consumer.New(cfg, service, producers, collector)
}


