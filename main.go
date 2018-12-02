package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func main() {
	err := InitDB("journal.sqlite")
	if err != nil {
		fmt.Println("Could not set up database")
		fmt.Println(err)
		return
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	Route(r)
	err = http.ListenAndServe(":3000", r)
	fmt.Println(err)
}
