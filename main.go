package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	err := InitDB("journal.sqlite")
	if err != nil {
		fmt.Println("Could not set up database")
		fmt.Println(err)
		return
	}
	err = prepUploadFolder()
	if err != nil {
		fmt.Println("Could not allocate upload folder!")
		fmt.Println(err)
		return
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// TODO: Customize the Recoverer to show a custom 500 page that has the same style as the rest of the app.
	r.Use(Recoverer)
	Route(r)
	err = http.ListenAndServe(":3000", r)
	fmt.Println(err)
}
