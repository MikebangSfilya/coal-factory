package equipment

import (
	"coalFactory/config"
	"errors"
	"log"
	"strings"
)

var (
	PickCost    = 5000
	VentCost    = 15000
	TrolleyCost = 50000
)

func Init(cfg *config.Configuration) {
	PickCost = cfg.PickCost
	VentCost = cfg.VentCost
	TrolleyCost = cfg.TrolleyCost
}

const (
	PickType    = "pick"
	VentType    = "vent"
	TrolleyType = "trolley"
)

type Equipments struct {
	Pick    Pick
	Vent    Vent
	Trolley Trolley
	AllBuy  bool
}

func NewEquipment() *Equipments {

	pick := Pick{Cost: PickCost}
	vent := Vent{Cost: VentCost}
	trolley := Trolley{Cost: TrolleyCost}
	return &Equipments{
		Pick:    pick,
		Vent:    vent,
		Trolley: trolley,
	}
}

func (e *Equipments) Buy(itemType string) (string, error) {
	item := strings.ToLower(itemType)

	switch item {
	case PickType:
		e.Pick.Buy()
		return PickType, nil
	case VentType:
		e.Vent.Buy()
		return VentType, nil
	case TrolleyType:
		e.Trolley.Buy()
		return TrolleyType, nil
	default:
		log.Print("not buy")
		return "", errors.New("")
	}
}

func (e *Equipments) AllBought() bool {
	if e.Pick.IsBought && e.Trolley.IsBought && e.Vent.IsBought {
		e.AllBuy = true
		return e.AllBuy
	}
	return false
}

type Pick struct {
	IsBought bool
	Cost     int
}

func (e *Pick) Buy() {
	e.IsBought = true
}

type Vent struct {
	IsBought bool
	Cost     int
}

func (e *Vent) Buy() {
	e.IsBought = true
}

type Trolley struct {
	IsBought bool
	Cost     int
}

func (e *Trolley) Buy() {
	e.IsBought = true
}
