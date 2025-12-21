package statistic

import (
	"coalFactory/equipment"
	"log/slog"
	"sync/atomic"
	"time"
)

type CompanyStats struct {
	Balance             *atomic.Int64
	TotalBalance        *atomic.Int64
	Income              *atomic.Int64
	Equipment           *equipment.Equipments
	Win                 bool
	TimeStarted         time.Time
	TimeEnd             *time.Time
	LittleMinersHired   int
	NormalMinersHired   int
	PowerfulMinersHired int
}

func New(equip *equipment.Equipments) *CompanyStats {

	totalBalance := &atomic.Int64{}
	Balance := &atomic.Int64{}
	income := &atomic.Int64{}

	return &CompanyStats{
		Balance:             Balance,
		TotalBalance:        totalBalance,
		Income:              income,
		Equipment:           equip,
		Win:                 false,
		TimeStarted:         time.Now(),
		TimeEnd:             nil,
		LittleMinersHired:   0,
		NormalMinersHired:   0,
		PowerfulMinersHired: 0,
	}
}

// Проверяет выигрышь если все куплено. Если все куплено отмечает время победы
func (cs *CompanyStats) CheckWinGame() (bool, error) {
	if !cs.Equipment.AllBuyed() {
		slog.Info("The user checked the victory, the victory was not achieved.", "INFO", errNotWin)
		return false, errNotWin
	}
	ended := time.Now()

	cs.TimeEnd = &ended
	cs.Win = true
	slog.Info("The user checked the victory, the victory was achieved.")
	return true, nil
}

func (cs *CompanyStats) CheckEquipment() *equipment.Equipments {
	return cs.Equipment
}

func (cs *CompanyStats) TimeCompleted() string {
	if cs.TimeEnd == nil {
		return ""
	}

	return cs.TimeEnd.Sub(cs.TimeStarted).String()
}

func (cs *CompanyStats) GetBalance() int64 {
	return cs.Balance.Load()
}

func (cs *CompanyStats) GetTotalBalance() int64 {
	return cs.TotalBalance.Load()
}

func (cs *CompanyStats) GetLittleMiners() int {
	return cs.LittleMinersHired
}

func (cs *CompanyStats) GetNormalMiners() int {
	return cs.NormalMinersHired
}

func (cs *CompanyStats) GetPowerfulMiners() int {
	return cs.PowerfulMinersHired
}
