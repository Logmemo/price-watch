package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ViewData struct {
	Title    string
	Heading1 string
	Products []string
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", basicHandler)

	server := &http.Server{
		Addr:    ":8011",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen to server", err)
	}

}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	tempData := ViewData{
		Title:    "Test main page",
		Heading1: "Overview",
	}
	tmpl, _ := template.ParseFiles("web/templates/main.html")
	tmpl.Execute(w, tempData)
}
