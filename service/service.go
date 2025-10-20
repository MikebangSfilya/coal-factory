package service

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/miners"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	GetMiners() map[uuid.UUID]factory.Miners
	GetMiner(id string) (factory.Miners, error)
	HireMiner(minerType miners.MinerType) (factory.Miners, error)

	GetBalance() int
	GetEq() equipment.Equipments

	WinGame() error
}

type EquipmentRepo interface {
	Buy(itemType string) (string, error)
}

type GameService struct {
	comp      CompanyRepo
	equipment EquipmentRepo
}

func New(company CompanyRepo, equipment EquipmentRepo) *GameService {
	return &GameService{
		comp:      company,
		equipment: equipment,
	}
}

func (gs *GameService) GetMiners() map[uuid.UUID]factory.Miners {
	miners := gs.comp.GetMiners()
	return miners
}

func (gs *GameService) GetMiner(id string) (factory.Miners, error) {
	miner, err := gs.comp.GetMiner(id)
	if err != nil {
		return nil, err
	}
	return miner, nil
}

func (gs *GameService) Hire(minerType miners.MinerType) (factory.Miners, error) {

	miner, err := gs.comp.HireMiner(minerType)
	if err != nil {
		return nil, err
	}

	return miner, nil
}

func (gs *GameService) Balance() int {
	return gs.comp.GetBalance()
}

func (gs *GameService) CheckWinGame() (bool, error) {
	err := gs.comp.WinGame()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (gs *GameService) Buy(item string) {
	gs.equipment.Buy(item)
}

func (gs *GameService) Items() equipment.Equipments {
	b := gs.comp.GetEq()
	return b
}
