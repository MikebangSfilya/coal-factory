package miners

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	MinerTypeNormal = "normal"
	NormalSalary    = 50
)

type NormalMiner struct {
	Id            uuid.UUID
	SleepDuration time.Duration
	CoalIncome    *atomic.Int64
	Energy        *atomic.Int64
}

func NewNormalMiner() *NormalMiner {
	const (
		timeSleep    = 2 * time.Second
		coal         = 3
		energyNormal = 45
	)

	energy := &atomic.Int64{}
	coalPower := &atomic.Int64{}
	energy.Add(energyNormal)
	coalPower.Add(coal)

	return &NormalMiner{
		Id:            uuid.New(),
		SleepDuration: timeSleep,
		CoalIncome:    coalPower,
		Energy:        energy,
	}
}

func (m *NormalMiner) Run(ctx context.Context) <-chan Coal {
	//Запускаем нашу горутину в бесконечно цикле
	// 1. Добываем coal
	// 2. отнимает от нашей energy 1
	// 3. Если энергия равна 0 то прекращаем работу
	// 4. Если энергия еще есть Засыпаем на 3 секунды
	// 5. Передаем по каналу количество добытого угля

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
		}

	}()

	return transferPoint

}

func (m *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		ID:        m.Id,
		MinerType: MinerTypeNormal,
		CoalPower: Coal(m.CoalIncome.Load()),
		Energy:    m.Energy.Load(),
		Cost:      NormalSalary,
	}
}
