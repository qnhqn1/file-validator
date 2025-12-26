package bootstrap

import (
	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/cache"
)


func InitCache(cfg *config.Config) (cache.Cache, func(), error) {
	c := cache.NewRedisCache(cfg.Redis)
	closer := func() { c.Close() }
	return c, closer, nil
}


