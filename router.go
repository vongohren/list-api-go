package main

import (
  "net/http"
  "github.com/Snorlock/shoppingApi/handlers"
  "github.com/Snorlock/shoppingApi/middleware"
  "github.com/Snorlock/shoppingApi/db"
  "github.com/Snorlock/mux"
)

type Router interface {
  GetRoute() []*mux.Route
}

func NewRouter(env *db.Env) *http.ServeMux {
  var authorize = true
  router := mux.NewRouter();
  router.Handle("/auth/{provider}", middleware.AuthHandler{middleware.Handler{env, !authorize}, handlers.BeginAuthHandler}).Methods("GET")
  router.Handle("/auth/{provider}/callback", middleware.AuthHandler{middleware.Handler{env, !authorize}, handlers.CallBack}).MakePrivate()
  router.Handle("/add/item", middleware.TokenHandler{middleware.Handler{env, authorize}, handlers.AddItemHandler}).Methods("POST")
  router.Handle("/add/list", middleware.TokenHandler{middleware.Handler{env, authorize}, handlers.AddListHandler}).Methods("POST")
  router.Handle("/lists", middleware.TokenHandler{middleware.Handler{env, authorize}, handlers.GetListsHandler}).Methods("GET", "OPTIONS")
  router.Handle("/lists/{id}", middleware.TokenHandler{middleware.Handler{env, authorize}, handlers.GetListDetailHandler}).Methods("GET", "OPTIONS")


  var routes = router.GetRoutes()
  router.Handle("/", middleware.IndexHandler{middleware.Handler{env, authorize}, handlers.IndexHandler, routes}).Methods("GET")


  mx := http.NewServeMux()
  mx.Handle("/", router)
  return mx
}
