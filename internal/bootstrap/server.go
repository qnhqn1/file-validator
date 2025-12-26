package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qnhqn1/file-validator/config"
	filevalidatorapi "github.com/qnhqn1/file-validator/internal/api/file_validator_api"
	"github.com/qnhqn1/file-validator/internal/consumer"
)


func AppRun(ctx context.Context, cfg *config.Config, api *filevalidatorapi.API, consumers *consumer.Manager) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: api.Router(),
	}

	go func() {
		log.Printf("%s слушает на :%d", cfg.ServiceName, cfg.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("ошибка HTTP сервера: %v", err)
		}
	}()

	consumerErrs := make(chan error, 1)
	go func() {
		consumerErrs <- consumers.Run(ctx)
	}()

	select {
	case <-ctx.Done():
	case err := <-consumerErrs:
		if err != nil {
			log.Printf("Kafka: ошибка консьюмера: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = server.Shutdown(shutdownCtx)
	consumers.Close()
	return nil
}


