package miners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinersInfo(t *testing.T) {
	type infoer interface {
		Info() MinerInfo
	}

	testCases := []struct {
		name          string
		createMiner   func() infoer
		wantType      MinerType
		wantCoalPower Coal
		wantEnergy    Energy
		wantCost      Coal
	}{
		{
			name:          "little miner",
			createMiner:   func() infoer { return NewLittleMiner() },
			wantType:      MinerTypeLittle,
			wantCoalPower: 1,
			wantEnergy:    30,
			wantCost:      LittleSalary,
		},
		{
			name:          "normal miner",
			createMiner:   func() infoer { return NewNormalMiner() },
			wantType:      MinerTypeNormal,
			wantCoalPower: 3,
			wantEnergy:    45,
			wantCost:      NormalSalary,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			miner := tc.createMiner()
			info := miner.Info()

			assert.Equal(t, tc.wantType, info.MinerType)
			assert.Equal(t, tc.wantCoalPower, info.CoalPower)
			assert.Equal(t, tc.wantEnergy, info.Energy)
			assert.Equal(t, tc.wantCost, info.Cost)
		})
	}
}

func TestLittleMiner_Run(t *testing.T) {
	t.Run("should stop when energy is depleted", func(t *testing.T) {
		miner := NewLittleMiner()
		miner.SleepDuration = time.Millisecond * 10
		miner.Energy = 2

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		coalChan := miner.Run(ctx)

		count := 0
		for range coalChan {
			count++
		}

		assert.Equal(t, 2, count, "Майнер должен был добыть ровно 2 угля")
		assert.Equal(t, Energy(0), miner.Energy, "Энергия должна быть на нуле")

	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		miner := NewLittleMiner()
		miner.SleepDuration = time.Hour

		ctx, cancel := context.WithCancel(context.Background())

		coalChan := miner.Run(ctx)

		cancel()

		select {
		case _, ok := <-coalChan:
			assert.False(t, ok, "Канал должен быть закрыт")
		case <-time.After(time.Millisecond * 100):
			t.Fatal("Горутина не завершилась вовремя после отмены контекста")
		}
	})
}

func TestNormalMiner_Run(t *testing.T) {
	t.Run("should stop when energy is depleted", func(t *testing.T) {
		miner := NewNormalMiner()
		miner.SleepDuration = time.Millisecond * 10
		miner.Energy = 2

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		coalChan := miner.Run(ctx)

		count := 0
		for v := range coalChan {
			count += int(v)
		}

		assert.Equal(t, 6, count, "Майнер должен был добыть ровно 6 угля")
		assert.Equal(t, Energy(0), miner.Energy, "Энергия должна быть на нуле")

	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		miner := NewLittleMiner()
		miner.SleepDuration = time.Hour

		ctx, cancel := context.WithCancel(context.Background())

		coalChan := miner.Run(ctx)

		cancel()

		select {
		case _, ok := <-coalChan:
			assert.False(t, ok, "Канал должен быть закрыт")
		case <-time.After(time.Millisecond * 100):
			t.Fatal("Горутина не завершилась вовремя после отмены контекста")
		}
	})
}

func TestPowerfulMiner_Run(t *testing.T) {
	t.Run("should stop when energy is depleted", func(t *testing.T) {
		miner := NewPowerfulMiner()
		miner.SleepDuration = time.Millisecond * 10
		miner.Energy = 2

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		coalChan := miner.Run(ctx)

		count := 0
		for v := range coalChan {
			count += int(v)
		}

		assert.Equal(t, 23, count, "Майнер должен был добыть ровно 6 угля")
		assert.Equal(t, Energy(0), miner.Energy, "Энергия должна быть на нуле")

	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		miner := NewLittleMiner()
		miner.SleepDuration = time.Hour

		ctx, cancel := context.WithCancel(context.Background())

		coalChan := miner.Run(ctx)

		cancel()

		select {
		case _, ok := <-coalChan:
			assert.False(t, ok, "Канал должен быть закрыт")
		case <-time.After(time.Millisecond * 100):
			t.Fatal("Горутина не завершилась вовремя после отмены контекста")
		}
	})
}
