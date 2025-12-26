package bootstrap

import (
	"github.com/qnhqn1/file-validator/internal/cache"
	"github.com/qnhqn1/file-validator/internal/services/validator"
	"github.com/qnhqn1/file-validator/internal/storage/pgstorage"
)


func InitValidatorService(storage *pgstorage.Storage, cache cache.Cache) validator.Service {
	return validator.New(storage, cache)
}


