package router

import (
	"goTest/internal/infrastructure/component"
	"goTest/internal/modules"
	"net/http"

	"github.com/go-chi/chi"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/ws", controllers.Messangerer.WebSocketHandler)

	return r
}
