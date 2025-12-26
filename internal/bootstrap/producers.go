package bootstrap

import (
	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/producer"
)


func InitProducers(cfg *config.Config) *producer.Manager {
	return producer.New(cfg)
}


