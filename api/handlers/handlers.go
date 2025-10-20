package handlers

import (
	dto "coalFactory/api/DTO"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/miners"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type HandleRepo interface {
	GetMiners() map[uuid.UUID]factory.Miners
	GetMiner(id string) (factory.Miners, error)
	Hire(minerType miners.MinerType) (factory.Miners, error)
	Balance() int
	CheckWinGame() (bool, error)
	Buy(item string)
	Items() equipment.Equipments
}

type Handlers struct {
	service HandleRepo
}

func New(handl HandleRepo) *Handlers {
	return &Handlers{
		service: handl,
	}
}

func (h *Handlers) Hire(w http.ResponseWriter, r *http.Request) {

	var dtoin dto.DTOHireMiner
	if err := json.NewDecoder(r.Body).Decode(&dtoin); err != nil {
		slog.Error("failed to decode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
		return
	}

	miner, err := h.service.Hire(miners.MinerType(dtoin.MinerType))
	if err != nil {
		slog.Error("Not enough coal for buy miner", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
	}

}

func (h *Handlers) GetMiners(w http.ResponseWriter, r *http.Request) {
	b := h.service.GetMiners()

	if err := json.NewEncoder(w).Encode(b); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
	}

}

func (h *Handlers) GetInfoMiner(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	miner, err := h.service.GetMiner(id)
	if err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
	}

}

func (h *Handlers) GetBal(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(h.service.Balance()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
	}

}

func (h *Handlers) CheckWin(w http.ResponseWriter, r *http.Request) {

	b, err := h.service.CheckWinGame()
	if err != nil {
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusPreconditionFailed)
		return
	}

	if b {
		w.Write([]byte("win"))
	}
}

// QueryParams
func (h *Handlers) BuyItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	eqName := query.Get("item")

	h.service.Buy(eqName)

}

func (h *Handlers) ItemsInfo(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(h.service.Items()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO, http.StatusInternalServerError)
	}
}
