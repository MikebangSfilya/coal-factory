package miners

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypeLittle = "little"
	LittleSalary    = 5
)

type LittleMiner struct {
	sync.RWMutex
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    Coal
	Energy        Energy
}

func NewLittleMiner() *LittleMiner {

	const (
		timeSleep    = 3 * time.Second
		coal         = 1
		energyLittle = 30
	)

	return &LittleMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coal,
		Energy:        energyLittle,
	}
}

func (m *LittleMiner) Run(ctx context.Context) <-chan Coal {

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
					slog.Info("Miner dead")
					return
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

		}

	}()

	return transferPoint
}

func (m *LittleMiner) Info() MinerInfo {
	m.RLock()
	defer m.RUnlock()

	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypeLittle,
		Energy:    m.Energy,
		CoalPower: m.CoalIncome,
		Cost:      LittleSalary,
	}
}
