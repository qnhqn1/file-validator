package filevalidatorapi

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/qnhqn1/file-validator/internal/api/swagger"
	"github.com/qnhqn1/file-validator/internal/metrics"
	"github.com/qnhqn1/file-validator/internal/services/validator"
)


type API struct {
	service       validator.Service
	serviceName   string
	enableSwagger bool
	collector     *metrics.Collector
	once          sync.Once
	swaggerSpec   []byte
}


func New(service validator.Service, serviceName string, collector *metrics.Collector) *API {
	return &API{
		service:       service,
		serviceName:   serviceName,
		enableSwagger: true, // assuming enableSwagger is true
		collector:     collector,
	}
}


func (a *API) Router() http.Handler {
	router := chi.NewRouter()
	router.Get("/health", a.health)
	router.Get("/metrics", a.collector.Handler())
	if a.enableSwagger {
		router.Get("/swagger", a.swaggerUI)
		router.Get("/swagger/validator.swagger.json", a.swaggerSpecHandler)
	}
	return router
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	body := map[string]string{
		"service": a.serviceName,
		"status":  "ok",
	}
	writeJSON(w, http.StatusOK, body)
}

func (a *API) swaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.25.0/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.25.0/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/swagger/validator.swagger.json',
            dom_id: '#swagger-ui',
        });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (a *API) swaggerSpecHandler(w http.ResponseWriter, r *http.Request) {
	a.once.Do(func() {
		a.swaggerSpec = swagger.Validator()
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(a.swaggerSpec)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

}


