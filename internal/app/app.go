package app

import (
	"context"
	"fmt"

	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/bootstrap"
	"github.com/qnhqn1/file-validator/internal/metrics"
)


func Run(ctx context.Context, cfg *config.Config) error {
	storage, err := bootstrap.InitPGStorage(ctx, cfg)
	if err != nil {
		return fmt.Errorf("инициализация хранилища: %w", err)
	}
	defer storage.Close()

	cache, closeCache, err := bootstrap.InitCache(cfg)
	if err != nil {
		return fmt.Errorf("инициализация кеша: %w", err)
	}
	defer closeCache()

	collector, err := metrics.New()
	if err != nil {
		return fmt.Errorf("инициализация метрик: %w", err)
	}

	service := bootstrap.InitValidatorService(storage, cache)
	api := bootstrap.InitValidatorAPI(service, cfg.ServiceName, collector)
	producers := bootstrap.InitProducers(cfg)
	consumers := bootstrap.InitConsumers(cfg, service, collector, producers)

	return bootstrap.AppRun(ctx, cfg, api, consumers)
}


