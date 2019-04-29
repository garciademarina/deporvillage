package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/garciademarina/deporvillage/pkg/adding"
	"github.com/garciademarina/deporvillage/pkg/handler"
	"github.com/garciademarina/deporvillage/pkg/listing"
	"github.com/garciademarina/deporvillage/pkg/updating"
)

// Server serves http requests.
type Server struct {
	Config
	Logger *log.Logger
	Router chi.Router
}

// Config represents server configuration
type Config struct {
	Port int
	Env  string
}

// NewConfig creates new config.
func NewConfig(port int, env string) Config {
	return Config{
		Port: port,
		Env:  env,
	}
}

// NewServer creates new Server.
func NewServer(config Config, logger *log.Logger, adder adding.Service, lister listing.Service, updater updating.Service) *Server {
	r := chi.NewRouter()

	registerRoutes(r, config, logger, lister, adder, updater)

	return &Server{
		Config: config,
		Logger: logger,
		Router: r,
	}
}

func registerRoutes(r chi.Router, config Config, logger *log.Logger, lister listing.Service, adder adding.Service, updater updating.Service) {
	r.Use(middleware.RequestID)

	handler := handler.NewOrderHandler(lister, adder, updater)

	r.Get("/order/{ID}", handler.GetOrder(logger))
	r.Post("/order", handler.AddOrder(logger))
	r.Put("/order", handler.UpdateStatusOrder(logger))
}

// Run runs the server.
func (s *Server) Run(ctx context.Context) error {
	s.Logger.Printf("Server started. PORT:%d\n", s.Port)

	http.Handle("/", s.Router)

	return http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), s.Router)
}

// ServeHTTP serve just one request.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
