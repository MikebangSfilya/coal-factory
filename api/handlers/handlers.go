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
	GetMiners(ctx context.Context) map[uuid.UUID]factory.Miners
	GetMiner(ctx context.Context, id string) (factory.Miners, error)
	Hire(ctx context.Context, minerType miners.MinerType) (factory.Miners, error)
	Balance(ctx context.Context) int
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

func New(handle HandleRepo) *Handlers {
	return &Handlers{
		service: handle,
	}
}

// Hire возвращает обработчик для найма минеров
func (h *Handlers) Hire() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), HireTimeOut)
		defer cancel()

		var dtoIn dto_in.DTOHireMiner
		if err := json.NewDecoder(r.Body).Decode(&dtoIn); err != nil {
			slog.Error("failed to decode JSON", "layer", "h", "operation", "decode", "error", err)
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}

		if err := dtoIn.Validate(); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, dto_in.ErrEmptyMinerType) {
				slog.Error("Empty type miner text", "error", err)
				status = http.StatusBadRequest
			} else if errors.Is(err, dto_in.ErrUnknownCommandMiner) {
				slog.Error("Unknown type miner text", "error", err)
				status = http.StatusUnprocessableEntity
			}

			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), status)
			return
		}

		miner, err := h.service.Hire(ctx, miners.MinerType(dtoIn.MinerType))
		if err != nil {
			slog.Error("Not enough coal for hire miner", "layer", "h", "operation", "hire", "error", err)
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusPaymentRequired)
			return
		}

		if err := json.NewEncoder(w).Encode(miner.Info()); err != nil {
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
		}
	}
}

// GetMiners возвращает обработчик списка всех минеров
func (h *Handlers) GetMiners() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
		defer cancel()

		minersMap := h.service.GetMiners(ctx)
		result := make(map[string]miners.MinerInfo, len(minersMap))
		for id, miner := range minersMap {
			result[id.String()] = miner.Info()
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
		}
	}
}

// GetInfoMiner возвращает обработчик инфо конкретного минера
func (h *Handlers) GetInfoMiner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
		}
	}
}

// GetBal возвращает обработчик баланса
func (h *Handlers) GetBal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
		defer cancel()

		if err := json.NewEncoder(w).Encode(h.service.Balance(ctx)); err != nil {
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
		}
	}
}

// CheckWin возвращает обработчик проверки победы и закрытия сервера
func (h *Handlers) CheckWin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
		defer cancel()

		stats, err := h.service.CheckWinGame(ctx)
		if err != nil {
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusPreconditionFailed)
			return
		}

		if err := json.NewEncoder(w).Encode(dto_out.NewDTOStats(stats)); err != nil {
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
			return
		}

		go func() {
			if h.serverClose != nil {
				if err := h.serverClose(); err != nil {
					slog.Debug("server close error", "error", err)
				}
			}
		}()
	}
}

// BuyItem возвращает обработчик покупки предмета
func (h *Handlers) BuyItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), NormalTimeOut)
		defer cancel()

		itemType := chi.URLParam(r, "type")
		_, err := h.service.Buy(ctx, itemType)
		if err != nil {
			errDTO := dto_out.NewErrorDTO(err)
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		}

		slog.Info("Item purchased successfully", "itemType", itemType)

		if err := json.NewEncoder(w).Encode(dto_out.NewDTORespItem(itemType)); err != nil {
			slog.Error("failed to encode JSON", "layer", "h", "operation", "encode", "error", err)
		}
	}
}
