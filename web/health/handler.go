package health

import (
	"net/http"

	"github.com/faceit/test/logger"
	healthcheck "github.com/faceit/test/services/health"
	"github.com/faceit/test/web"
	"github.com/faceit/test/web/middleware"
	"github.com/gorilla/mux"
)

// Handler is a web events handler struct
type Handler struct {
	router     *mux.Router
	log        logger.Logger
	middleware middleware.Middleware
	check      *healthcheck.Check
}

// NewHandler creates new healtcheck handler instancce
func NewHandler(router *mux.Router, l logger.Logger, m middleware.Middleware, c *healthcheck.Check) {
	h := Handler{
		router:     router,
		log:        l,
		middleware: m,
		check:      c,
	}

	apiV1 := router.PathPrefix("/v1").Subrouter()

	apiV1.HandleFunc("/health", h.middleware.SetContextHeader(http.HandlerFunc(h.Check))).
		Methods(http.MethodGet)
}

// Check handles Get health check request
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	newHealth(web.NewResponse(w, h.log), h.check).Do(web.NewRequest(r))
}
