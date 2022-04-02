package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/executive", handlers.Repo.Executive)
	mux.Get("/general", handlers.Repo.General)

	mux.Get("/search-ava", handlers.Repo.SearchAva)
	mux.Post("/search-ava", handlers.Repo.PostSearchAva)
	mux.Post("/search-ava-json", handlers.Repo.AvailabilityJson)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/res-summary", handlers.Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
