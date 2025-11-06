package handlers

import (
	"coalFactory/api/DTO/dto_in"
	"coalFactory/api/DTO/dto_out"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type HandleRepo interface {
	//miners
	GetMiners(ctx context.Context) map[uuid.UUID]factory.Miners
	GetMiner(ctx context.Context, id string) (factory.Miners, error)
	Hire(ctx context.Context, minerType miners.MinerType) (factory.Miners, error)
	//stats
	Balance(ctx context.Context) int
	//items
	Buy(ctx context.Context, item string) (*equipment.Equipments, error)
	Items(ctx context.Context) equipment.Equipments
	CheckWinGame(ctx context.Context) (statistic.CompanyStats, error)
}

const (
	HireTimeOut   = 30 * time.Second
	NormalTimeOut = 10 * time.Second
)

type Handlers struct {
	service     HandleRepo
	serverClose func() error
}

func New(handl HandleRepo) *Handlers {
	return &Handlers{
		service: handl,
	}
}

func (handlers *Handlers) CloseServer(f func() error) {
	handlers.serverClose = f
}

func (h *Handlers) Hire(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), HireTimeOut)
	defer cancel()

	var dtoin dto_in.DTOHireMiner
	if err := json.NewDecoder(r.Body).Decode(&dtoin); err != nil {
		slog.Error(
			"failed to decode JSON",
			"layer", "handlers",
			"operation", "decode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if err := dtoin.Validate(); err != nil {
		if errors.Is(err, dto_in.ErrEmptyMinerType) {
			slog.Error("Empty type miner text", "error", err)
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		} else if errors.Is(err, dto_in.ErrUnknownCommandMiner) {
			slog.Error("Unknown  type miner text", "error", err)
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusUnprocessableEntity)
			return
		} else {
			slog.Error("Internal server error", "error", err)
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	miner, err := h.service.Hire(ctx, miners.MinerType(dtoin.MinerType))
	if err != nil {
		slog.Error(
			"Not enough coal for hire miner",
			"layer", "handlers",
			"operation", "hire",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusPaymentRequired)
		return
	}

	if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetMiners(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	minersMap := h.service.GetMiners(ctx)

	result := make(map[string]miners.MinerInfo, len(minersMap)) // TODO перенести в сервисный слой обработку
	for id, miner := range minersMap {
		result[id.String()] = miner.Info()
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetInfoMiner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	id := chi.URLParam(r, "id")

	miner, err := h.service.GetMiner(ctx, id)
	if err != nil {
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetBal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	if err := json.NewEncoder(w).Encode(h.service.Balance(ctx)); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) CheckWin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	stats, err := h.service.CheckWinGame(ctx)
	if err != nil {
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusPreconditionFailed)
		return
	}

	dtoStats := dto_out.NewDTOStats(stats)

	if err := json.NewEncoder(w).Encode(dtoStats); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	go func() {
		if err := h.serverClose(); err != nil {
			slog.Debug("server close error", "error", err)
		}
	}()
}

// QueryParams либо в JSON теле
func (h *Handlers) BuyItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	itemType := chi.URLParam(r, "type")

	_, err := h.service.Buy(ctx, itemType)
	if err != nil {
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	slog.Info("Item purchased successfully",
		"layer", "handlers",
		"operation", "buy",
		"itemType", itemType)

	dtoResp := dto_out.NewDTORespItem(itemType)

	if err := json.NewEncoder(w).Encode(dtoResp); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) ItemsInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
	defer cancel()

	items := h.service.Items(ctx)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error(
			"failed to encode JSON",
			"layer", "handlers",
			"operation", "encode",
			"error", err,
		)
		errDTO := dto_out.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}
}
