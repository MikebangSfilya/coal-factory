package statistic

import (
	"coalFactory/equipment"
	"log/slog"
	"sync/atomic"
	"time"
)

type CompanyStats struct {
	Balance     *atomic.Int64
	income      *atomic.Int64 //На самом деле пока никак не работает, надо подумать как аккамулировать весь доход с каналов
	Equipmet    *equipment.Equipments
	Win         bool
	TimeStarted time.Time
	TimeEnd     *time.Time
}

func New(equip *equipment.Equipments) *CompanyStats {

	income := &atomic.Int64{}
	Balance := &atomic.Int64{}

	return &CompanyStats{
		Balance:     Balance,
		income:      income,
		Equipmet:    equip,
		Win:         false,
		TimeStarted: time.Now(),
	}
}

// Проверяет выигрышь если все куплено. Если все куплено отмечает время победы
func (cs *CompanyStats) CheckWinGame() (bool, error) {
	if !cs.Equipmet.AllBuyed() {
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
	return cs.Equipmet
}
