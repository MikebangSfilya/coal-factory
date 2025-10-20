package server

import (
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
}

type Server struct {
	handlers HTTPRepo
}

func New(handlers HTTPRepo) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Post("/miners", s.handlers.Hire)
	r.Get("/miners/{id}", s.handlers.GetInfoMiner)
	r.Get("/miners", s.handlers.GetMiners)

	r.Get("/miners/balance", s.handlers.GetBal)
	r.Get("/miners/win", s.handlers.CheckWin)

	r.Post("/miners/items", s.handlers.BuyItem)
	r.Get("/miners/items", s.handlers.ItemsInfo)

	fmt.Print("started")
	return http.ListenAndServe(":9091", r)
}
