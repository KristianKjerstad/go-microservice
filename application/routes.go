package application

import (
	"net/http"
	"project/handlers"
	"project/repository/order"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader((http.StatusOK))
		writer.Write([]byte("No content"))
	})

	router.Route("/orders", a.loadOrderRoutes)

	a.router = router

}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handlers.Order{
		Repo: &order.RedisRepo{
			Client: a.redisDB,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{}", orderHandler.DeleteByID)
}
