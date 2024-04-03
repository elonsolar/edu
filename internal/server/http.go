package server

import (
	"context"
	adminV1 "edu/api/admin/v1"
	assistant "edu/api/assistant/v1"
	courseV1 "edu/api/course/v1"
	customerV1 "edu/api/customer/v1"
	productV1 "edu/api/product/v1"
	"edu/internal/conf"
	"edu/internal/domain/model"
	"edu/internal/service"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt4 "github.com/golang-jwt/jwt/v4"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/api.admin.v1.Auth/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	dossier *service.DossierService,
	customer *service.CustomerService,
	auth *service.AuthService,
	product *service.ProductService,

	task *service.TaskService,
	course *service.CourseService,
	logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			selector.Server(
				jwt.Server(func(token *jwt4.Token) (interface{}, error) {
					return []byte("edu"), nil
				}, jwt.WithSigningMethod(jwt4.SigningMethodHS256), jwt.WithClaims(func() jwt4.Claims {
					return &model.Claims{}
				})),
			).Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	// if c.Http.Timeout != nil {
	// 	opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	// }
	opts = append(opts, http.Timeout(time.Second*40))
	srv := http.NewServer(opts...)
	adminV1.RegisterDossierHTTPServer(srv, dossier)
	adminV1.RegisterAuthHTTPServer(srv, auth)

	customerV1.RegisterCustomerHTTPServer(srv, customer)

	courseV1.RegisterCourseHTTPServer(srv, course)
	assistant.RegisterTaskHTTPServer(srv, task)
	productV1.RegisterProductHTTPServer(srv, product)

	return srv
}
