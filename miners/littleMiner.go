package miners

import (
	"context"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypeLittle = "little"
	LittleSalary    = 5
)

type LittleMiner struct {
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    *atomic.Int64
	Energy        *atomic.Int64
}

func NewLittleMiner() *LittleMiner {

	const (
		timeSleep    = 3 * time.Second
		coal         = 1
		energyLittle = 30
	)
	energy := &atomic.Int64{}
	coalPower := &atomic.Int64{}
	energy.Add(energyLittle)
	coalPower.Add(coal)

	return &LittleMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coalPower,
		Energy:        energy,
	}
}

func (m *LittleMiner) Run(ctx context.Context) <-chan Coal {
	// #TODO добавить передачу дохода в статистику

	transferPoint := make(chan Coal)

	go func() {
		defer close(transferPoint)

		for m.Energy.Load() > 0 {
			select {
			case <-ctx.Done():
				slog.Info("Contex stopped", "contex", ctx)
				return
			case <-time.After(m.SleepDuration):

			}
			select {
			case <-ctx.Done():
				return
			case transferPoint <- Coal(m.CoalIncome.Load()):
				m.Energy.Add(-1)
			}
		}

		if m.Energy.Load() == 0 {
			log.Printf("Работник %d окончил", m.Id)
		}

	}()

	return transferPoint
}

func (m *LittleMiner) Info() MinerInfo {
	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypeLittle,
		Energy:    m.Energy.Load(),
		CoalPower: Coal(m.CoalIncome.Load()),
		Cost:      LittleSalary,
	}
}
