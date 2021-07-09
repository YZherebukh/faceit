package country

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/faceit/test/logger"
	"github.com/faceit/test/services/country"
	"github.com/faceit/test/web"
	"github.com/faceit/test/web/middleware"
)

// Handler is a web events handler struct
type Handler struct {
	router     *mux.Router
	log        logger.Logger
	middleware middleware.Middleware
	country    *country.Country
}

// NewHandler creates new user handler instancce
func NewHandler(router *mux.Router, l logger.Logger, m middleware.Middleware, c *country.Country) {
	h := Handler{
		router:     router,
		log:        l,
		middleware: m,
		country:    c,
	}

	apiV1 := router.PathPrefix("/v1").Subrouter()

	apiV1.HandleFunc("/countries", h.middleware.SetContextHeader(http.HandlerFunc(h.All))).
		Methods(http.MethodGet)
}

// All handles Get All countries requests
func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	newAll(web.NewResponse(w, h.log), h.country).Do(web.NewRequest(r))
}
