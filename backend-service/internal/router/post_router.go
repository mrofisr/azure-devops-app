package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrofisr/azure-devops/internal/handler"
)

func PostRouter(postHandler *handler.PostHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", postHandler.FindAll)
	r.Get("/{id}", postHandler.FindByID)
	r.Get("/count", postHandler.Count)
	r.Post("/", postHandler.Create)
	r.Put("/", postHandler.Update)
	r.Delete("/{id}", postHandler.Delete)
	return r
}
