package miners

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypePowerful = "powerful"
	PowerfulSalary    = 450
)

type PowerfulMiner struct {
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    *atomic.Int64
	Energy        *atomic.Int64
}

func NewPowerfulMiner() *PowerfulMiner {
	const (
		timeSleep      = 1 * time.Second
		coal           = 10
		energyPowerful = 60
	)

	energy := &atomic.Int64{}
	coalPower := &atomic.Int64{}
	energy.Add(energyPowerful)
	coalPower.Add(coal)

	return &PowerfulMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coalPower,
		Energy:        energy,
	}
}

func (m *PowerfulMiner) Run(ctx context.Context) <-chan Coal {

	const increaseIncome = 3

	transferPoint := make(chan Coal)

	go func() {
		defer close(transferPoint)

		for m.Energy.Load() > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.SleepDuration):
			}
			select {
			case <-ctx.Done():
				return
			case transferPoint <- Coal(m.CoalIncome.Load()):
				m.Energy.Add(-1)
			}
			m.CoalIncome.Add(increaseIncome)
		}

	}()

	return transferPoint

}

func (m *PowerfulMiner) Info() MinerInfo {
	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypePowerful,
		CoalPower: Coal(m.CoalIncome.Load()),
		Energy:    m.Energy.Load(),
		Cost:      PowerfulSalary,
	}
}
