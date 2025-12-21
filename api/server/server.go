package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPRepo interface {
	Hire(w http.ResponseWriter, r *http.Request)
	GetInfoMiner(w http.ResponseWriter, r *http.Request)
	GetMiners(w http.ResponseWriter, r *http.Request)

	GetBal(w http.ResponseWriter, r *http.Request)
	CheckWin(w http.ResponseWriter, r *http.Request)

	BuyItem(w http.ResponseWriter, r *http.Request)
	ItemsInfo(w http.ResponseWriter, r *http.Request)

	CloseServer(f func() error)
}

type Server struct {
	handlers HTTPRepo
	server   *http.Server
}

func New(handlers HTTPRepo) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	r.Route("/miners", func(r chi.Router) {
		r.Post("/", s.handlers.Hire)
		r.Get("/{id}", s.handlers.GetInfoMiner)
		r.Get("/", s.handlers.GetMiners)
	})
	r.Route("/items", func(r chi.Router) {
		r.Post("/{type}", s.handlers.BuyItem)
		r.Get("/", s.handlers.ItemsInfo)
	})

	r.Get("/balance", s.handlers.GetBal)
	r.Get("/win", s.handlers.CheckWin)

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	s.handlers.CloseServer(s.server.Close)
	fmt.Println("started")
	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
