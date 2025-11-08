package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (s *HealthController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.LogDebug("health check ok")
		w.Write([]byte("ok"))
	})
	return r
}
