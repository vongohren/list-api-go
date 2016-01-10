package middleware

import (
    "net/http"
    "log"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
    error
    Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
    Code int
    Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
    return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
    return se.Code
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.

func (h TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  _ = "breakpoint"
  var i HandlerMethods
  i = h
  handleCors(w, r)
  // Stop here if its Preflighted OPTIONS request
  if r.Method == "OPTIONS" {
      return
  }
  id, failed := checkAuthorization(i, w, r)
  if !failed {
    err := h.HandleWithToken(h.Env, id, w, r)
    if err != nil {
      handleError(err, w)
    }
  }
}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var i HandlerMethods
  i = h
  _, failed := checkAuthorization(i, w, r)
  if !failed {
    err := h.HandleWithRoutes(h.Env, h.Routes, w, r)
    if err != nil {
      handleError(err, w)
    }
  }
}

func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err := h.Handle(h.Env, w, r)
  if err != nil {
    handleError(err, w)
  }
}

func handleError(err interface{}, w http.ResponseWriter) {
  switch e := err.(type) {
  case Error:
      // We can retrieve the status here and write out a specific
      // HTTP status code.
      log.Printf("HTTP %d - %s", e.Status(), e)
      http.Error(w, e.Error(), e.Status())
  default:
      // Any error types we don't specifically look out for default
      // to serving a HTTP 500
      log.Printf("Default handler for error %s", e)
      http.Error(w, http.StatusText(http.StatusInternalServerError),
          http.StatusInternalServerError)
  }
}

func handleCors(w http.ResponseWriter, r *http.Request) {
  if origin := r.Header.Get("Origin"); origin != "" {
      w.Header().Set("Access-Control-Allow-Origin", origin)
      w.Header().Set("Access-Control-Allow-Credentials", "true")
      w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
      w.Header().Set("Access-Control-Allow-Headers",
          "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
  }
}
