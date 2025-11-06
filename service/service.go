package service

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"context"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	GetMiners(ctx context.Context) map[uuid.UUID]factory.Miners
	GetMiner(ctx context.Context, id string) (factory.Miners, error)
	HireMiner(ctx context.Context, minerType miners.MinerType) (factory.Miners, error)

	GetBalance(ctx context.Context) int
	GetEq(ctx context.Context) equipment.Equipments

	WinGame(ctx context.Context) (statistic.CompanyStats, error)
	Buy(ctx context.Context, itemType string) (*equipment.Equipments, error)
}

type GameService struct {
	comp CompanyRepo
}

func New(company CompanyRepo) *GameService {
	return &GameService{
		comp: company,
	}
}

func (gs *GameService) GetMiners(ctx context.Context) map[uuid.UUID]factory.Miners {
	return gs.comp.GetMiners(ctx)
}

func (gs *GameService) GetMiner(ctx context.Context, id string) (factory.Miners, error) {
	return gs.comp.GetMiner(ctx, id)
}

func (gs *GameService) Hire(ctx context.Context, minerType miners.MinerType) (factory.Miners, error) {
	return gs.comp.HireMiner(ctx, minerType)
}

func (gs *GameService) Balance(ctx context.Context) int {
	return gs.comp.GetBalance(ctx)
}

func (gs *GameService) CheckWinGame(ctx context.Context) (statistic.CompanyStats, error) {
	return gs.comp.WinGame(ctx)
}

func (gs *GameService) Buy(ctx context.Context, item string) (*equipment.Equipments, error) {
	return gs.comp.Buy(ctx, item)
}

func (gs *GameService) Items(ctx context.Context) equipment.Equipments {
	return gs.comp.GetEq(ctx)
}
