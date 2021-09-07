package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Port         string
	Router       *chi.Mux
	RouteHandler *RouteHandler
}

func NewServer(port string) *Server {
	server := Server{
		Port:         port,
		Router:       chi.NewRouter(),
		RouteHandler: &RouteHandler{},
	}

	return &server
}

func (server *Server) Run() {
	server.Router = chi.NewRouter()
	server.Router.Use(middleware.Logger)

	server.registerRoutes()

	http.ListenAndServe(":"+server.Port, server.Router)
}

func (server *Server) registerRoutes() {
	server.Router.Post("/register", server.RouteHandler.Register())
	server.Router.Get("/get", server.RouteHandler.Register())
	server.Router.Delete("/delete", server.RouteHandler.Register())
}
