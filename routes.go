package main

import (
	"MyWeb/pkg/configs"
	"MyWeb/pkg/handlers"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routesController(app *configs.AppConfig) http.Handler {
	// we will mux  == HTTP request multiplexer

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(10 * time.Second))
	// mux.Use(middleware.L
	mux.Use(LogRequestInfo)

	mux.Use(NoSurf)
	mux.Use(SetupSession)

	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)

	mux.Get("/login", handlers.Repo.LoginHandler)
	mux.Post("/login", handlers.Repo.PostLoginHandler)

	mux.Get("/logout", handlers.Repo.LogoutHandler)

	mux.Get("/page", handlers.Repo.PageHandler)
	mux.Get("/makepost", handlers.Repo.MakePostHandler)

	mux.Post("/makepost", handlers.Repo.PostMakePostHandler)

	mux.Get("/article-received", handlers.Repo.ArticleReceived)

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
