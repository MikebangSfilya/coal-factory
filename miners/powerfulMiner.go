package miners

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypePowerful = "powerful"
	PowerfulSalary    = 450
)

type PowerfulMiner struct {
	sync.RWMutex
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    Coal
	Energy        Energy
}

func NewPowerfulMiner() *PowerfulMiner {
	const (
		timeSleep      = 1 * time.Second
		coal           = 10
		energyPowerful = 60
	)

	return &PowerfulMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coal,
		Energy:        energyPowerful,
	}
}

func (m *PowerfulMiner) Run(ctx context.Context) <-chan Coal {

	const increaseIncome = 3

	transferPoint := make(chan Coal)
	currInc := m.CoalIncome

	go func() {
		defer close(transferPoint)

		ticker := time.NewTicker(m.SleepDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.Lock()
				if m.Energy <= 0 {
					m.Unlock()
					slog.Debug("Miner pow dead")
					return
				}
			}

			m.Energy--
			send := currInc
			currInc += increaseIncome
			m.CoalIncome = send

			m.Unlock()

			select {
			case <-ctx.Done():
				return
			case transferPoint <- send:
			}

		}

	}()

	return transferPoint

}

func (m *PowerfulMiner) Info() MinerInfo {
	m.RLock()
	defer m.RUnlock()
	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypePowerful,
		CoalPower: m.CoalIncome,
		Energy:    m.Energy,
		Cost:      PowerfulSalary,
	}
}
