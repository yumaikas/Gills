package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Route(r chi.Router) {
	r.Get("/", Home)
}

// TODO: This is going to provide HTTP handlers that are routed by main.go
// And then it's going to load state from the database and feed it into
// the various routers.
func Home(w http.ResponseWriter, r *http.Request) {
	TempHomeView(w)
}
