package service

import (
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewDossierService, NewAuthService, NewTaskService, NewCustomerService, NewCourseService, NewProductService)

var (
	tracer = otel.Tracer("service")
)
