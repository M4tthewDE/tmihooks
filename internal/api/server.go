package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/tmi"
)

type Server struct {
	Port         string
	server       *http.Server
	Router       *chi.Mux
	RouteHandler *RouteHandler
	stopChan     chan os.Signal
}

func NewServer(config *config.Config, reader *tmi.Reader) *Server {
	server := Server{
		Port:         config.Server.Port,
		server:       &http.Server{Addr: ":" + config.Server.Port},
		Router:       chi.NewRouter(),
		RouteHandler: NewRouteHandler(config, reader),
		stopChan:     make(chan os.Signal, 1),
	}
	signal.Notify(server.stopChan, os.Interrupt)

	return &server
}

func (server *Server) Run() {
	server.Router = chi.NewRouter()
	server.Router.Use(middleware.Logger)

	server.registerRoutes()
	server.server.Handler = server.Router

	go func() {
		err := server.server.ListenAndServe()
		if err != nil {
			log.Println("server closed")
		}
	}()

	<-server.stopChan
	log.Println("stopping main server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("done")
}

func (server *Server) registerRoutes() {
	server.Router.Post("/register", server.RouteHandler.Register())
	server.Router.Get("/get", server.RouteHandler.Get())
	server.Router.Delete("/delete", server.RouteHandler.Delete())
	server.Router.Post("/shutdown", server.RouteHandler.Shutdown())
}
