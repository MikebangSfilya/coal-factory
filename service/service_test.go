package service_test

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"coalFactory/service"
	"context"
	"testing"

	"github.com/google/uuid"
)

type MockMiners struct {
	RunFunc  func(ctx context.Context) <-chan miners.Coal
	InfoFunc func() miners.MinerInfo
}

func (m *MockMiners) Run(ctx context.Context) <-chan miners.Coal {
	if m.RunFunc != nil {
		return m.RunFunc(ctx)
	}
	ch := make(chan miners.Coal)
	close(ch)
	return ch
}

func (m *MockMiners) Info() miners.MinerInfo {
	if m.InfoFunc != nil {
		return m.InfoFunc()
	}
	return miners.MinerInfo{}
}

type MockCompanyRepo struct {
	GetMinersFunc func() map[uuid.UUID]factory.Miners
	GetMinerFunc  func(id string) (factory.Miners, error)
	HireMinerFunc func(minerType miners.MinerType) (factory.Miners, error)

	GetBalanceFunc func() int
	GetEqFunc      func() equipment.Equipments

	WinGameFunc func() (statistic.CompanyStats, error)
	BuyFunc     func(itemType string) (*equipment.Equipments, error)
}

func (m *MockCompanyRepo) GetMiners() map[uuid.UUID]factory.Miners {
	return map[uuid.UUID]factory.Miners{}
}

func (m *MockCompanyRepo) GetMiner(id string) (factory.Miners, error) {
	return nil, nil
}

func (m *MockCompanyRepo) HireMiner(minerType miners.MinerType) (factory.Miners, error) {
	if m.HireMinerFunc != nil {
		return m.HireMinerFunc(minerType)
	}
	return nil, factory.ErrNotEnoughMoney
}

func (m *MockCompanyRepo) GetBalance() int {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc()
	}
	return 0
}

func (m *MockCompanyRepo) GetEq() equipment.Equipments {
	return equipment.Equipments{}
}

func (m *MockCompanyRepo) WinGame() (statistic.CompanyStats, error) {
	return statistic.CompanyStats{}, nil
}

func (m *MockCompanyRepo) Buy(itemType string) (*equipment.Equipments, error) {
	return nil, nil
}

func Test_HireMiner(t *testing.T) {

	testCases := []struct {
		name          string
		minerType     miners.MinerType
		balance       int
		wantMinerType miners.MinerType
	}{
		{
			name:          "hire little miner",
			minerType:     miners.MinerTypeLittle,
			balance:       1000,
			wantMinerType: miners.MinerTypeLittle,
		},
		{
			name:          "hire normal miner",
			minerType:     miners.MinerTypeNormal,
			balance:       1000,
			wantMinerType: miners.MinerTypeNormal,
		},
		{
			name:          "hire powerful miner",
			minerType:     miners.MinerTypePowerful,
			balance:       1000,
			wantMinerType: miners.MinerTypePowerful,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			comp := &MockCompanyRepo{}

			comp.HireMinerFunc = func(minerType miners.MinerType) (factory.Miners, error) {
				return &MockMiners{
					InfoFunc: func() miners.MinerInfo {
						return miners.MinerInfo{
							ID:        uuid.New(),
							MinerType: minerType,
							CoalPower: 1,
							Energy:    30,
							Cost:      5,
						}
					},
				}, nil
			}

			comp.GetBalanceFunc = func() int {
				return tc.balance
			}

			service := service.New(comp)

			miner, err := service.Hire(tc.minerType)

			if err != nil {
				t.Fatalf("Ожидался успешный найм, но ошибка %v", err)
			}

			if miner == nil {
				t.Fatalf("Ожидалось наличие майнера, но получили nil")
			}

			if miner.Info().MinerType != tc.wantMinerType {
				t.Errorf("Ожидалось что майнер %s, но получили %s", tc.wantMinerType, miner.Info().MinerType)
			}

			t.Log("Проверки прошли успешно")

		})
	}

}
