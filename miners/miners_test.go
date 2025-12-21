package miners_test

import (
	"coalFactory/factory"
	"coalFactory/miners"
	"testing"
)

func TestMinerInfo(t *testing.T) {
	testCases := []struct {
		name          string
		createMiner   func() factory.Miners
		wantType      miners.MinerType
		wantCoalPower miners.Coal
		wantEnergy    int64
		wantCost      int
	}{
		{
			name: "little miner info",
			createMiner: func() factory.Miners {
				return miners.NewLittleMiner()
			},
			wantType:      miners.MinerTypeLittle,
			wantCoalPower: 1,
			wantEnergy:    30,
			wantCost:      miners.LittleSalary,
		},
		{
			name: "normal miner info",
			createMiner: func() factory.Miners {
				return miners.NewNormalMiner()
			},
			wantType:      miners.MinerTypeNormal,
			wantCoalPower: 3,
			wantEnergy:    45,
			wantCost:      miners.NormalSalary,
		},
		{
			name: "powerful miner info",
			createMiner: func() factory.Miners {
				return miners.NewPowerfulMiner()
			},
			wantType:      miners.MinerTypePowerful,
			wantCoalPower: 10,
			wantEnergy:    60,
			wantCost:      miners.PowerfulSalary,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			miner := tc.createMiner()
			info := miner.Info()

			if info.MinerType != tc.wantType {
				t.Errorf("Want type miner %s", tc.wantType)
			}
			if info.CoalPower != tc.wantCoalPower {
				t.Errorf("Mining power: get: %d, expected: %d", info.CoalPower, tc.wantCoalPower)
			}
			if info.Energy != tc.wantEnergy {
				t.Errorf("Energy: get: %d, expected: %d", info.Energy, tc.wantEnergy)
			}
			if info.Cost != int64(tc.wantCost) {
				t.Errorf("Cost: get: %d, expected:  %d", info.Cost, tc.wantCost)
			}
		})
	}
}
