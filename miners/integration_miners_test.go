package miners_test

import (
	"coalFactory/factory"
	"coalFactory/miners"
	"context"
	"testing"
	"time"
)

func Test_PowerfulMiner_Integration(t *testing.T) {
	miner := miners.NewPowerfulMiner()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := miner.Run(ctx)

	var coals []miners.Coal

	for i := 0; i < 3; i++ {
		select {
		case coal := <-ch:
			coals = append(coals, coal)
			t.Logf("Coal collection iteration %d: %d coal collected", i+1, coal)
		case <-time.After(5 * time.Second):
			t.Errorf("Coal has not been collected for a long time")
		}
	}

	expected := []miners.Coal{10, 13, 16}
	for i := 0; i < len(coals); i++ {
		if coals[i] != expected[i] {
			t.Errorf("Coal amount mismatch, expected %d, got %d", expected[i], coals[i])
		}
	}

}

func TestRunMiners_Integration(t *testing.T) {

	deltaStart := 1 * time.Second
	testCases := []struct {
		name      string
		miner     factory.Miners
		coal      int64
		timeSleep time.Duration
	}{
		{
			name:      "Testing little miner",
			miner:     miners.NewLittleMiner(),
			coal:      1,
			timeSleep: 3*time.Second + deltaStart,
		},
		{
			name:      "Testing normal miner",
			miner:     miners.NewNormalMiner(),
			coal:      3,
			timeSleep: 2*time.Second + deltaStart,
		},
		{
			name:      "Testing powerful miner",
			miner:     miners.NewPowerfulMiner(),
			coal:      10,
			timeSleep: 1*time.Second + deltaStart,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			testChan := tc.miner.Run(ctx)
			select {
			case coal := <-testChan:
				if coal != miners.Coal(tc.coal) {
					t.Errorf("Expected to receive %d, but got %d", tc.coal, coal)
				}
			case <-time.After(tc.timeSleep + 1*time.Second):
				t.Errorf("Miner is taking too long to produce coal")
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
					t.Errorf("Channel should have been closed after context cancellation")
				}
			case <-time.After(1 * time.Second):
				t.Errorf("Channel was not closed within 1 second")
			}
		})
	}

}
