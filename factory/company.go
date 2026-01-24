package factory

import (
	"coalFactory/equipment"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	pick    = "pick"
	vent    = "vent"
	trolley = "trolley"
)

type Miners interface {
	Run(ctx context.Context) <-chan miners.Coal
	Info() miners.MinerInfo
}

type Company struct {
	Miners map[uuid.UUID]Miners
	Income chan miners.Coal //Канал для передачи наших денег в баланс

	CompanyContext context.Context
	CompanyStop    context.CancelFunc

	mu sync.RWMutex

	Stats *statistic.CompanyStats
}

func NewCompany(ctx context.Context, equip *equipment.Equipments) *Company {
	companyContext, companyStop := context.WithCancel(ctx)
	slog.Info( //Лог на будущее если будет возможность делать несколько игроков
		"create new Company",
		"layer", "company",
		"operation", "create",
	)
	return &Company{
		Miners:         make(map[uuid.UUID]Miners),
		Income:         make(chan miners.Coal),
		CompanyContext: companyContext,
		CompanyStop:    companyStop,
		Stats:          statistic.New(equip),
	}
}

func Start(ctx context.Context, equip *equipment.Equipments) *Company {
	comp := NewCompany(ctx, equip)
	go comp.PassiveIncome()
	go comp.RaiseBalance()

	slog.Info("start new Company")
	return comp
}

func (c *Company) GetMiners(ctx context.Context) map[uuid.UUID]Miners {
	copyMap := make(map[uuid.UUID]Miners, len(c.Miners))
	for k, v := range c.Miners {
		copyMap[k] = v
	}
	return copyMap
}

func (c *Company) GetMiner(ctx context.Context, id string) (Miners, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	miner, ok := c.Miners[uuid]
	if !ok {
		return nil, ErrMinerNotExist
	}
	return miner, nil
}

func (c *Company) HireMiner(ctx context.Context, minerType miners.MinerType) (Miners, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var miner Miners

	switch minerType {
	case miners.MinerTypeLittle:
		if c.Stats.Balance.Load() >= miners.LittleSalary {
			miner = miners.NewLittleMiner()
			c.Stats.Balance.Add(-miners.LittleSalary)
			c.Stats.LittleMinersHired++
		} else {
			return nil, ErrNotEnoughMoney
		}
	case miners.MinerTypeNormal:
		if c.Stats.Balance.Load() >= miners.NormalSalary {
			miner = miners.NewNormalMiner()
			c.Stats.Balance.Add(-miners.NormalSalary)
			c.Stats.NormalMinersHired++
		} else {
			return nil, ErrNotEnoughMoney
		}
	case miners.MinerTypePowerful:
		if c.Stats.Balance.Load() >= miners.PowerfulSalary {
			miner = miners.NewPowerfulMiner()
			c.Stats.Balance.Add(-miners.PowerfulSalary)
			c.Stats.PowerfulMinersHired++
		} else {
			return nil, ErrNotEnoughMoney
		}
	}

	c.Miners[miner.Info().ID] = miner

	c.StartMiner(miner)

	return miner, nil
}

func (c *Company) StartMiner(miner Miners) {
	coalTransferPoint := miner.Run(c.CompanyContext)
	go func() {
		slog.Info(
			"start miner job",
			"layer", "company",
			"operation", "run",
			"miner", miner.Info(),
			"cost", miner.Info().Cost,
		)

		c.Stats.Income.Add(int64(miner.Info().CoalPower))
		defer func() {
			c.Stats.Income.Add(-int64(miner.Info().CoalPower))
			slog.Info(
				"miner finished work",
				"layer", "company",
				"operation", "stop",
				"miner_id", miner.Info().ID,
				"reason", "energy depleted",
			)
		}()

		for val := range coalTransferPoint {
			c.Income <- val
		}
		c.mu.Lock()
		delete(c.Miners, miner.Info().ID)
		c.mu.Unlock()
	}()
}

func (c *Company) RaiseBalance() {

	go func() {
		for {
			select {
			case <-c.CompanyContext.Done():
				return
			case val := <-c.Income:
				c.Stats.Balance.Add(int64(val))
				c.Stats.TotalBalance.Add(int64(val))
			}
		}
	}()

}

// Запуск нашего пассивного дохода равного 1 единице
func (c *Company) PassiveIncome() {
	c.Stats.Income.Add(1)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.CompanyContext.Done():
			return
		case <-ticker.C:
			select {
			case <-c.CompanyContext.Done():
				return
			case c.Income <- 1:

			}
		}
	}
}

func (c *Company) GetBalance(ctx context.Context) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return int(c.Stats.Balance.Load())
}

func (c *Company) WinGame(ctx context.Context) (statistic.CompanyStats, error) {
	win, err := c.Stats.CheckWinGame()
	if err != nil {
		return statistic.CompanyStats{}, err
	}
	if win {
		c.CompanyStop()
	}

	return *c.Stats, nil
}

func (c *Company) GetEq(ctx context.Context) equipment.Equipments {
	return *c.Stats.CheckEquipment()
}

func (c *Company) Buy(ctx context.Context, itemType string) (*equipment.Equipments, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	itemTypeLower := strings.ToLower(itemType)
	switch itemTypeLower {
	case pick:
		if c.Stats.Balance.Load() >= int64(equipment.PickCost) {
			slog.Info(
				"equipment purchased",
				"layer", "company",
				"operation", "buy",
				"item", pick,
				"cost", equipment.PickCost,
			)
			c.Stats.Equipment.Buy(pick)
			c.Stats.Balance.Add(-int64(equipment.PickCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	case vent:
		if c.Stats.Balance.Load() >= int64(equipment.VentCost) {
			slog.Info(
				"equipment purchased",
				"layer", "company",
				"operation", "buy",
				"item", vent,
				"cost", equipment.VentCost,
			)
			c.Stats.Equipment.Buy(vent)
			c.Stats.Balance.Add(-int64(equipment.VentCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	case trolley:
		if c.Stats.Balance.Load() >= int64(equipment.TrolleyCost) {
			slog.Info(
				"equipment purchased",
				"layer", "company",
				"operation", "buy",
				"item", trolley,
				"cost", equipment.TrolleyCost,
			)
			c.Stats.Equipment.Buy(trolley)
			c.Stats.Balance.Add(-int64(equipment.TrolleyCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	default:
		slog.Warn(
			"unknown itemType",
			"layer", "company",
			"operation", "buy",
			"item", "unknown",
			"cost", "-1",
		)
		return nil, errors.New("unknown item type: " + itemType)
	}
	return c.Stats.Equipment, nil
}

func (c *Company) SetBalance(balance int) {
	c.Stats.Balance.Store(int64(balance))
}
