package controllers

import "github.com/go-chi/chi/v5"

type IController interface {
	MapController() *chi.Mux
}
