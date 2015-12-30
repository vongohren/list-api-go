package middleware

import (
    "net/http"
    "github.com/Snorlock/shoppingApi/db"
    "github.com/Snorlock/mux"
)

type Handler struct {
  *db.Env
  Authorize bool
}

type IndexHandler struct {
  Handler
  HandleWithRoutes func(e *db.Env, routes []*mux.Route, w http.ResponseWriter, r *http.Request) error
  Routes []*mux.Route
}

type TokenHandler struct {
  Handler
  HandleWithToken func(e *db.Env, token interface{}, w http.ResponseWriter, r *http.Request) error
}

type AuthHandler struct {
  Handler
  Handle func(e *db.Env, w http.ResponseWriter, r *http.Request) error
}

type HandlerMethods interface {
  shallAuthorize() bool
}

func (h Handler) shallAuthorize() bool {
    return h.Authorize
}
