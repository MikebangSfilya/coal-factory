package handlers

import (
	dto "coalFactory/api/DTO"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/miners"
	"encoding/json"
	"errors"
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
	Buy(item string) (*equipment.Equipments, error)
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
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if err := dtoin.Validate(); err != nil {
		if errors.Is(err, dto.ErrEmptyMinerType) {
			slog.Error("Empty type miner text", "error", err)
			errDTO := dto.NewErrorDto(err)
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		} else if errors.Is(err, dto.ErrEmptyMinerType) {
			slog.Error("Unknow type miner text", "error", err)
			errDTO := dto.NewErrorDto(err)
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		} else {
			slog.Error("InternalServerError", "error", err)
			errDTO := dto.NewErrorDto(err)
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	miner, err := h.service.Hire(miners.MinerType(dtoin.MinerType))
	if err != nil {
		slog.Error("Not enough coal for buy miner", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetMiners(w http.ResponseWriter, r *http.Request) {
	b := h.service.GetMiners()

	if err := json.NewEncoder(w).Encode(b); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
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
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetBal(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(h.service.Balance()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) CheckWin(w http.ResponseWriter, r *http.Request) {

	b, err := h.service.CheckWinGame()
	if err != nil {
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusPreconditionFailed)
		return
	}

	if b {
		w.Write([]byte("win"))
	} else {
		w.Write([]byte("not win yet"))
	}
}

// QueryParams либо в JSON файле
func (h *Handlers) BuyItem(w http.ResponseWriter, r *http.Request) {

	var itemType string
	var dtoItem dto.DTORBuyItem
	if err := json.NewDecoder(r.Body).Decode(&dtoItem); err == nil {
		if err := dtoItem.Validate(); err != nil {
			if errors.Is(err, dto.ErrEmptyItemType) {
				slog.Error("The user sent an empty json", "error", err)
				errDTO := dto.NewErrorDto(err)
				http.Error(w, errDTO.ToString(), http.StatusBadRequest)
				return
			} else if errors.Is(err, dto.ErrUnknowCommandItem) {
				slog.Error("The user sent wrong itemType", "error", err)
				errDTO := dto.NewErrorDto(err)
				http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			} else {
				slog.Error("Internal Server Error", "error", err)
				errDTO := dto.NewErrorDto(err)
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}
		}
		itemType = dtoItem.ItemType
	} else {
		itemType = r.URL.Query().Get("item")
		if itemType == "" {
			slog.Error("No item specified in JSON or query")
			errDTO := dto.NewErrorDto(ErrCannotParse)
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		}
	}

	_, err := h.service.Buy(itemType)
	if err != nil {
		slog.Error("Cant buy")
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	dtoResp := dto.NewResp(itemType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dtoResp); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) ItemsInfo(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(h.service.Items()); err != nil {
		slog.Error("failed to encode json", "error", err)
		errDTO := dto.NewErrorDto(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}
}
