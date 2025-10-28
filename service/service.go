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
	GetMiners() map[uuid.UUID]factory.Miners
	GetMiner(id string) (factory.Miners, error)
	HireMiner(ctx context.Context, minerType miners.MinerType) (factory.Miners, error)

	GetBalance() int
	GetEq() equipment.Equipments

	WinGame() (statistic.CompanyStats, error)
	Buy(itemType string) (*equipment.Equipments, error)
}

type GameService struct {
	comp CompanyRepo
}

func New(company CompanyRepo) *GameService {
	return &GameService{
		comp: company,
	}
}

func (gs *GameService) GetMiners() map[uuid.UUID]factory.Miners {
	return gs.comp.GetMiners()
}

func (gs *GameService) GetMiner(id string) (factory.Miners, error) {
	return gs.comp.GetMiner(id)
}

func (gs *GameService) Hire(ctx context.Context, minerType miners.MinerType) (factory.Miners, error) {
	return gs.comp.HireMiner(ctx, minerType)
}

func (gs *GameService) Balance() int {
	return gs.comp.GetBalance()
}

func (gs *GameService) CheckWinGame() (statistic.CompanyStats, error) {
	return gs.comp.WinGame()
}

func (gs *GameService) Buy(item string) (*equipment.Equipments, error) {
	return gs.comp.Buy(item)
}

func (gs *GameService) Items() equipment.Equipments {
	return gs.comp.GetEq()
}
