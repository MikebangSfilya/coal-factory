package server

import (
	"context"
	"net/http"
	"time"

	_ "coalFactory/docs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/swaggo/http-swagger"
)

type HTTPRepo interface {
	Hire() http.HandlerFunc
	GetInfoMiner() http.HandlerFunc
	GetMiners() http.HandlerFunc
	GetBal() http.HandlerFunc
	CheckWin() http.HandlerFunc
	BuyItem() http.HandlerFunc
	ItemsInfo() http.HandlerFunc
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
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Группировка маршрутов
	r.Route("/miners", func(r chi.Router) {
		r.Post("/", handlers.Hire())
		r.Get("/{id}", handlers.GetInfoMiner())
		r.Get("/", handlers.GetMiners())
	})

	r.Route("/items", func(r chi.Router) {
		r.Post("/{type}", handlers.BuyItem())
		r.Get("/", handlers.ItemsInfo())
	})

	r.Get("/balance", handlers.GetBal())
	r.Get("/win", handlers.CheckWin())

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  5 * time.Second,
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
