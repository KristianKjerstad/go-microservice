package application

import (
	"net/http"
	"project/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader((http.StatusOK))
		writer.Write([]byte("No content"))
	})

	router.Route("/orders", loadOrderRoutes)

	return router

}

func loadOrderRoutes(router chi.Router) {
	orderHandler := &handlers.Order{}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{}", orderHandler.DeleteByID)
}
