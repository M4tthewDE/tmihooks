package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m4tthewde/tmihooks/internal/config"
)

type Server struct {
	Port         string
	Router       *chi.Mux
	RouteHandler *RouteHandler
}

func NewServer(config *config.Config) *Server {
	server := Server{
		Port:         config.Server.Port,
		Router:       chi.NewRouter(),
		RouteHandler: NewRouteHandler(config),
	}

	return &server
}

func (server *Server) Run() {
	server.Router = chi.NewRouter()
	server.Router.Use(middleware.Logger)

	server.registerRoutes()

	err := http.ListenAndServe(":"+server.Port, server.Router)
	if err != nil {
		panic(err)
	}
}

func (server *Server) registerRoutes() {
	server.Router.Post("/register", server.RouteHandler.Register())
	server.Router.Get("/get", server.RouteHandler.Register())
	server.Router.Delete("/delete", server.RouteHandler.Register())
}
