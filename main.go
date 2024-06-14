package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen to server")
	}
}

func basicHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world"))

}
