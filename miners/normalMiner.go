package miners

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypeNormal = "normal"
	NormalSalary    = 50
)

type NormalMiner struct {
	sync.RWMutex
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    Coal
	Energy        Energy
}

func NewNormalMiner() *NormalMiner {
	const (
		timeSleep    = 2 * time.Second
		coal         = 3
		energyNormal = 45
	)
	return &NormalMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coal,
		Energy:        energyNormal,
	}
}

func (m *NormalMiner) Run(ctx context.Context) <-chan Coal {

	transferPoint := make(chan Coal)

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
					slog.Info("Normal miner stopped")
					return
				}
			}

			income := m.CoalIncome
			m.Energy--
			m.Unlock()

			select {
			case <-ctx.Done():
				return
			case transferPoint <- income:
			}
		}

	}()

	return transferPoint

}

func (m *NormalMiner) Info() MinerInfo {
	m.RLock()
	defer m.RUnlock()
	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypeNormal,
		CoalPower: m.CoalIncome,
		Energy:    m.Energy,
		Cost:      NormalSalary,
	}
}
