package miners_test

import (
	"coalFactory/miners"
	"context"
	"testing"
	"time"
)

type Miner interface {
	Run(ctx context.Context) <-chan miners.Coal
	Info() miners.MinerInfo
}

func Test_PowerfulMiner(t *testing.T) {
	miner := miners.NewPowerfulMiner()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := miner.Run(ctx)

	var Coals []miners.Coal

	for i := 0; i < 3; i++ {
		select {
		case coal := <-ch:
			Coals = append(Coals, coal)
			t.Logf("Итерация сбора угля %d: %d угля собрано", i+1, coal)
		case <-time.After(5 * time.Second):
			t.Errorf("Уголь не собирается уже долгое время")
		}
	}

	expected := []miners.Coal{10, 13, 16}
	for i := 0; i < len(Coals); i++ {
		if Coals[i] != expected[i] {
			t.Errorf("Количество угля не совпадает, ожидалось %d, получено %d", expected[i], Coals[i])
		}
	}

}

func TestTestRunMiners(t *testing.T) {

	deltaStart := 1 * time.Second
	testCases := []struct {
		name      string
		miner     Miner
		coal      int64
		timeSleep time.Duration
		energy    int
	}{
		{
			name:      "Testing little miner",
			miner:     miners.NewLittleMiner(),
			coal:      1,
			timeSleep: 3*time.Second + deltaStart,
			energy:    1,
		},
		{
			name:      "Testing normal miner",
			miner:     miners.NewNormalMiner(),
			coal:      3,
			timeSleep: 2*time.Second + deltaStart,
			energy:    1,
		},
		{
			name:      "Testing powerful miner",
			miner:     miners.NewPowerfulMiner(),
			coal:      10,
			timeSleep: 1*time.Second + deltaStart,
			energy:    1,
		},
	}

	for _, tc := range testCases {

		t.Run(t.Name(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			testChan := tc.miner.Run(ctx)
			select {
			case coal := <-testChan:
				if coal != miners.Coal(tc.coal) {
					t.Errorf("Ожидалось получение %d, было получено %d", tc.coal, coal)
				}
			case <-time.After(tc.timeSleep + 1*time.Second):
				t.Errorf("Майнер не добывает уголь слишком долго")
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			ch := tc.miner.Run(ctx)

			<-ch

			cancel()

			select {
			case _, ok := <-ch:
				if ok {
					t.Errorf("канал должен был быть закрыт после отмены контекста")
				}
			case <-time.After(1 * time.Second):
				t.Errorf("канал не закрыт спустя 1 секунду")
			}
		})
	}

}
