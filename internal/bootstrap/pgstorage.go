package bootstrap

import (
	"context"

	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/storage/pgstorage"
	"github.com/qnhqn1/file-validator/internal/storage/sharding"
)


func InitPGStorage(ctx context.Context, cfg *config.Config) (*pgstorage.Storage, error) {
	manager, err := sharding.NewShardManager(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}
	return pgstorage.New(ctx, manager)
}


