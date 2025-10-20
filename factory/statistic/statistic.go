package statistic

import (
	"coalFactory/equipment"
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
		return false, errNotWin
	}
	ended := time.Now()

	cs.TimeEnd = &ended
	cs.Win = true
	return true, nil
}

func (cs *CompanyStats) CheckEquipment() *equipment.Equipments {
	return cs.Equipmet
}
