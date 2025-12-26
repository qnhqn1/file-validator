package bootstrap

import (
	filevalidatorapi "github.com/qnhqn1/file-validator/internal/api/file_validator_api"
	"github.com/qnhqn1/file-validator/internal/metrics"
	"github.com/qnhqn1/file-validator/internal/services/validator"
)


func InitValidatorAPI(service validator.Service, serviceName string, collector *metrics.Collector) *filevalidatorapi.API {
	return filevalidatorapi.New(service, serviceName, collector)
}


