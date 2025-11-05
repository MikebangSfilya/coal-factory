package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPRepo interface {
	//miners
	Hire(w http.ResponseWriter, r *http.Request)
	GetInfoMiner(w http.ResponseWriter, r *http.Request)
	GetMiners(w http.ResponseWriter, r *http.Request)
	//Information
	GetBal(w http.ResponseWriter, r *http.Request)
	CheckWin(w http.ResponseWriter, r *http.Request)
	//Items
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

	r.Post("/miners", s.handlers.Hire)
	r.Get("/miners/{id}", s.handlers.GetInfoMiner)
	r.Get("/miners", s.handlers.GetMiners)

	r.Get("/balance", s.handlers.GetBal)
	r.Get("/win", s.handlers.CheckWin)

	r.Post("/items/{type}", s.handlers.BuyItem)
	r.Get("/items", s.handlers.ItemsInfo)

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
