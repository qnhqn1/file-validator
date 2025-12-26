package main

import (
	"context"
	"log"
	"os"

	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/app"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	if err := app.Run(context.Background(), cfg); err != nil {
		log.Fatalf("приложение завершено: %v", err)
	}
}


