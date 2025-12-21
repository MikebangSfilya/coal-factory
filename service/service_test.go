package service_test

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"coalFactory/service"
	"context"
	"errors"
	"testing"
	"time"

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
	GetMinersFunc func(ctx context.Context) map[uuid.UUID]factory.Miners
	GetMinerFunc  func(ctx context.Context, id string) (factory.Miners, error)
	HireMinerFunc func(ctx context.Context, minerType miners.MinerType) (factory.Miners, error)

	GetBalanceFunc func(ctx context.Context) int
	GetEqFunc      func(ctx context.Context) equipment.Equipments

	WinGameFunc func(ctx context.Context) (statistic.CompanyStats, error)
	BuyFunc     func(ctx context.Context, itemType string) (*equipment.Equipments, error)
}

func (m *MockCompanyRepo) GetMiners(ctx context.Context) map[uuid.UUID]factory.Miners {
	return map[uuid.UUID]factory.Miners{}
}

func (m *MockCompanyRepo) GetMiner(ctx context.Context, id string) (factory.Miners, error) {
	return nil, nil
}

func (m *MockCompanyRepo) HireMiner(ctx context.Context, minerType miners.MinerType) (factory.Miners, error) {
	if m.HireMinerFunc != nil {
		return m.HireMinerFunc(ctx, minerType)
	}
	return nil, factory.ErrNotEnoughMoney
}

func (m *MockCompanyRepo) GetBalance(ctx context.Context) int {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc(ctx)
	}
	return 0
}

func (m *MockCompanyRepo) GetEq(ctx context.Context) equipment.Equipments {
	return equipment.Equipments{}
}

func (m *MockCompanyRepo) WinGame(ctx context.Context) (statistic.CompanyStats, error) {
	return statistic.CompanyStats{}, nil
}

func (m *MockCompanyRepo) Buy(ctx context.Context, itemType string) (*equipment.Equipments, error) {
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

			comp.HireMinerFunc = func(ctx context.Context, minerType miners.MinerType) (factory.Miners, error) {
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

			comp.GetBalanceFunc = func(ctx context.Context) int {
				return tc.balance
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			service := service.New(comp)

			miner, err := service.Hire(ctx, tc.minerType)

			if err != nil {
				t.Fatalf("Expected successful hiring, but got error: %v", err)
			}

			if miner == nil {
				t.Fatalf("Expected a miner instance, but got nil")
			}

			if miner.Info().MinerType != tc.wantMinerType {
				t.Errorf("Expected miner type %s, but got %s", tc.wantMinerType, miner.Info().MinerType)
			}

			t.Log("All checks passed successfully")

		})
	}

}

func TestHire_ContextTimeout(t *testing.T) {
	mockRepo := &MockCompanyRepo{
		HireMinerFunc: func(ctx context.Context, minerType miners.MinerType) (factory.Miners, error) {
			time.Sleep(2 * time.Second)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return &MockMiners{}, nil
			}
		},
	}

	service := service.New(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := service.Hire(ctx, miners.MinerTypeLittle)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected Deadline Exceeded, got: %v", err)
	}
}
