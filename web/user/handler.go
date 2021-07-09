package user

import (
	"net/http"

	"github.com/faceit/test/config"
	"github.com/faceit/test/logger"
	"github.com/faceit/test/queue"
	"github.com/faceit/test/services/hasher"
	"github.com/faceit/test/services/password"
	"github.com/faceit/test/services/user"
	"github.com/faceit/test/store"
	"github.com/faceit/test/web"
	"github.com/faceit/test/web/middleware"

	"github.com/gorilla/mux"
)

// Handler is a web events handler struct
type Handler struct {
	router      *mux.Router
	notifierCFG config.Notifier
	log         logger.Logger
	middleware  middleware.Middleware
	queue       queue.Queue
	user        *user.User
	country     *store.Country
	password    *password.Password
	hssher      *hasher.Hasher
}

// NewHandler creates new user handler instancce
func NewHandler(r *mux.Router, l logger.Logger, m middleware.Middleware,
	u *user.User, c *store.Country, p *password.Password, hash *hasher.Hasher, q queue.Queue) {
	h := Handler{
		router:     r,
		log:        l,
		middleware: m,
		queue:      q,
		user:       u,
		country:    c,
		password:   p,
		hssher:     hash,
	}

	apiV1 := h.router.PathPrefix("/v1").Subrouter()

	apiV1.HandleFunc("/user", h.middleware.SetContextHeader(http.HandlerFunc(h.All))).
		Methods(http.MethodGet)
	apiV1.HandleFunc("/user", h.middleware.SetContextHeader(http.HandlerFunc(h.Create))).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/user/{id}", h.middleware.SetContextHeader(http.HandlerFunc(h.One))).
		Methods(http.MethodGet)
	apiV1.HandleFunc("/user/{id}", h.middleware.SetContextHeader(http.HandlerFunc(h.Update))).
		Methods(http.MethodPut)
	apiV1.HandleFunc("/user/{id}/password", h.middleware.SetContextHeader(http.HandlerFunc(h.UpdatePassword))).
		Methods(http.MethodPut)
	apiV1.HandleFunc("/user/{id}", h.middleware.SetContextHeader(http.HandlerFunc(h.Delete))).
		Methods(http.MethodDelete)
}

// All handles Get All users requests
// it returns all users by filter, if one was provided
// otherwise, it's just returns all users
func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	newAll(web.NewResponse(w, h.log), h.user).Do(web.NewRequest(r))
}

// One handles Get One event requests by userID
func (h *Handler) One(w http.ResponseWriter, r *http.Request) {
	newOne(web.NewResponse(w, h.log), h.user).Do(web.NewRequest(r))
}

// Create handles POST Create user requests
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	newCreate(web.NewResponse(w, h.log), h.user, h.queue, h.notifierCFG.OnCreate()).Do(web.NewRequest(r))
}

// Update handles PUT Update user requests
// ID and Password can not be updated with this request
// since ID is autogenerated by database and
// password has it's own update call
// to perform update user must send his current password
// for authorization
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	newUpdate(web.NewResponse(w, h.log), h.user, h.queue, h.notifierCFG.OnUpdate()).Do(web.NewRequest(r))
}

// UpdatePassword handles PUT UpdatePassword user requests
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	newUpdatePassword(web.NewResponse(w, h.log), h.password).Do(web.NewRequest(r))
}

// Delete handles delete request
// it requires users password in order to delete user
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	newDelete(web.NewResponse(w, h.log), h.user, h.queue, h.notifierCFG.OnDelete()).Do(web.NewRequest(r))
}