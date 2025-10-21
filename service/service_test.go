package service_test

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/miners"
	"coalFactory/service"
	"context"
	"testing"
)

func TestGerMiners(t *testing.T) {
	equip := equipment.NewEquipmet()

	company := factory.NewCompany(context.Background(), equip)

	service := service.New(company)

	mapa := service.GetMiners()

	if mapa == nil {
		t.Fatal("отсутствует карта для хранения майнеров")
	}

}

func TestHire(t *testing.T) {

	testCases := []struct {
		name         string
		minerTypes   string
		balanceStart int
		balanceEnd   int
	}{
		{
			name:         "hire little miner",
			minerTypes:   "little",
			balanceStart: 500,
			balanceEnd:   495,
		},
		{
			name:         "hire normal miner",
			minerTypes:   "normal",
			balanceStart: 500,
			balanceEnd:   450,
		},
		{
			name:         "hire powerful miner",
			minerTypes:   "powerful",
			balanceStart: 500,
			balanceEnd:   50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eqip := equipment.NewEquipmet()
			comp := factory.NewCompany(context.Background(), eqip)
			serv := service.New(comp)

			comp.SetBalance(tc.balanceStart)

			res, err := serv.Hire(miners.MinerType(tc.minerTypes))
			if err != nil {
				t.Errorf("не получилось нанять майнера")
			}
			gotBalance := comp.GetBalance()

			if gotBalance != tc.balanceEnd {
				t.Errorf("баланс неверный, ожидалось %d, получилось %d", tc.balanceEnd, gotBalance)
			}
			if res == nil {
				t.Errorf("нанять майнера не вышло")
			}

		})
	}

}

func TestGetMiner(t *testing.T) {
	equip := equipment.NewEquipmet()

	company := factory.NewCompany(context.Background(), equip)

	service := service.New(company)

	company.SetBalance(1000)

	littleMiner, err := company.HireMiner("little")
	if err != nil {
		t.Errorf("не получилсоь нанять")
	}

	idMiner := littleMiner.Info().ID

	miner, err := service.GetMiner(idMiner.String())
	if err != nil {
		t.Errorf("Не удалось найти шахтера %v", err)
	}

	if miner == nil {
		t.Fatal("Должен был вернуться шахтер")
	}

	if miner.Info().ID != littleMiner.Info().ID {
		t.Errorf("метод отработал неверно, найден не тот")
	}

}
