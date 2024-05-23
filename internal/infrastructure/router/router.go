package router

import (
	"goTest/internal/infrastructure/component"
	"goTest/internal/modules"
	"goTest/internal/router"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter(controllers *modules.Controllers, components *component.Components) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", router.NewApiRouter(controllers, components))
	return r
}
