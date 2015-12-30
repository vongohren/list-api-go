package main

import(
  "net/http"
  "github.com/Snorlock/shoppingApi/db"
)

func NewServer(addr string) *http.Server {
  //Setup Database
  env := db.StartDatabase()

	// Setup router
	var router = NewRouter(env)

	// Create and start server
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
