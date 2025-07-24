package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oyinetare/go-docker-microservice/api"
	"github.com/oyinetare/go-docker-microservice/repository"
)

type Server struct {
	repo   repository.RepositoryInterface
	router *mux.Router
	port   int
}

func New(repo repository.RepositoryInterface, port int) *Server {
	return &Server{
		repo: repo,
		port: port,
	}
}

func (s *Server) Start() error {
	s.router = mux.NewRouter()

	// Add logging middleware
	s.router.Use(loggingMiddleware)

	// Create API handlers
	userAPI := api.NewUserAPI(s.repo)

	// Register routes
	s.router.HandleFunc("/users", userAPI.GetUsers).Methods("GET")
	s.router.HandleFunc("/search", userAPI.SearchUser).Methods("GET")

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Server starting on port %d", s.port)
	return http.ListenAndServe(addr, s.router)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
