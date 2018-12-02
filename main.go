package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logging)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	})
	http.ListenAndServe(":3000", r)
}
