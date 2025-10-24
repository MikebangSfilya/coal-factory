package factory

import (
	"coalFactory/equipment"
	"coalFactory/factory/statistic"
	"coalFactory/miners"
	"context"
	"errors"
	"log"
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

func Start(equip *equipment.Equipments) *Company {
	comp := NewCompany(context.Background(), equip)
	go comp.PassiveIncome()
	go comp.RaiseBalance()
	go comp.ShowIncome()
	slog.Info(
		"start company layer",
		"layer", "company",
		"operation", "run",
	)
	return comp
}

// Возвращает мапу наших рабочих
func (c *Company) GetMiners() map[uuid.UUID]Miners {
	copyMap := make(map[uuid.UUID]Miners, len(c.Miners))
	for k, v := range c.Miners {
		copyMap[k] = v
	}
	return copyMap
}

func (c *Company) GetMiner(id string) (Miners, error) {
	iduuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	miner, ok := c.Miners[iduuid]
	if !ok {
		return nil, ErrMinerMotExist
	}
	return miner, nil
}

func (c *Company) HireMiner(minerType miners.MinerType) (Miners, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var miner Miners

	switch minerType {
	case miners.MinerTypeLittle:
		if c.Stats.Balance.Load() >= miners.LittleSalary {
			miner = miners.NewLittleMiner()
			c.Stats.Balance.Add(-miners.LittleSalary)
		} else {
			return nil, ErrNotEnoughMoney
		}
	case miners.MinerTypeNormal:
		if c.Stats.Balance.Load() >= miners.NormalSalary {
			miner = miners.NewNormalMiner()
			c.Stats.Balance.Add(-miners.NormalSalary)
		} else {
			return nil, ErrNotEnoughMoney
		}
	case miners.MinerTypePowerful:
		if c.Stats.Balance.Load() >= miners.PowerfulSalary {
			miner = miners.NewPowerfulMiner()
			c.Stats.Balance.Add(-miners.PowerfulSalary)
		} else {
			return nil, ErrNotEnoughMoney
		}
	}

	c.Miners[miner.Info().ID] = miner

	c.StartMiner(miner)

	return miner, nil
}

func (c *Company) StartMiner(miner Miners) {
	coalTranserPoint := miner.Run(context.Background())
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

		for val := range coalTranserPoint {
			c.Income <- val
		}
		delete(c.Miners, miner.Info().ID)
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
				c.Stats.TotalBalanced.Add(int64(val))
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

func (c *Company) GetBalance() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return int(c.Stats.Balance.Load())
}

func (c *Company) ShowIncome() {
	for {
		c.mu.RLock()
		b := c.Stats.Income
		log.Println(b.Load())
		time.Sleep(1 * time.Second)
		c.mu.RUnlock()
	}
}

func (c *Company) WinGame() (statistic.CompanyStats, error) {
	win, err := c.Stats.CheckWinGame()
	if err != nil {
		return statistic.CompanyStats{}, err
	}
	if win {
		c.CompanyStop()
	}

	return *c.Stats, nil
}

func (c *Company) GetEq() equipment.Equipments {
	return *c.Stats.CheckEquipment()
}

func (c *Company) Buy(itemType string) (*equipment.Equipments, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	itemTypelow := strings.ToLower(itemType)
	switch itemTypelow {
	case pick:
		if c.Stats.Balance.Load() >= int64(equipment.PickCost) {
			slog.Info(
				"equipment purchared",
				"layer", "company",
				"operation", "buy",
				"item", pick,
				"cost", equipment.PickCost,
			)
			c.Stats.Equipmet.Buy(pick)
			c.Stats.Balance.Add(-int64(equipment.PickCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	case vent:
		if c.Stats.Balance.Load() >= int64(equipment.VentCost) {
			slog.Info(
				"equipment purchared",
				"layer", "company",
				"operation", "buy",
				"item", vent,
				"cost", equipment.VentCost,
			)
			c.Stats.Equipmet.Buy(vent)
			c.Stats.Balance.Add(-int64(equipment.VentCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	case trolley:
		if c.Stats.Balance.Load() >= int64(equipment.TrolleyCost) {
			slog.Info(
				"equipment purchared",
				"layer", "company",
				"operation", "buy",
				"item", trolley,
				"cost", equipment.TrolleyCost,
			)
			c.Stats.Equipmet.Buy(trolley)
			c.Stats.Balance.Add(-int64(equipment.TrolleyCost))
		} else {
			return nil, ErrNotEnoughMoney
		}
	default:
		slog.Warn(
			"Uncnown itemType",
			"layer", "company",
			"operation", "buy",
			"item", "unknow",
			"cost", "-1",
		)
		return nil, errors.New("unknown item type: " + itemType)
	}
	return c.Stats.Equipmet, nil
}

// method for unit testing
func (c *Company) SetBalance(balance int) {
	c.Stats.Balance.Store(int64(balance))
}
