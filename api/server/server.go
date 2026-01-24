package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type HTTPRepo interface {
	Hire(w http.ResponseWriter, r *http.Request)
	GetInfoMiner(w http.ResponseWriter, r *http.Request)
	GetMiners(w http.ResponseWriter, r *http.Request)
	GetBal(w http.ResponseWriter, r *http.Request)
	CheckWin(w http.ResponseWriter, r *http.Request)
	BuyItem(w http.ResponseWriter, r *http.Request)
	ItemsInfo(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	httpServer *http.Server
}

func New(addr string, handlers HTTPRepo) *Server {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Группировка маршрутов
	r.Route("/miners", func(r chi.Router) {
		r.Post("/", handlers.Hire)
		r.Get("/{id}", handlers.GetInfoMiner)
		r.Get("/", handlers.GetMiners)
	})

	r.Route("/items", func(r chi.Router) {
		r.Post("/{type}", handlers.BuyItem)
		r.Get("/", handlers.ItemsInfo)
	})

	r.Get("/balance", handlers.GetBal)
	r.Get("/win", handlers.CheckWin)

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  5 * time.Second, // Хорошая практика
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
