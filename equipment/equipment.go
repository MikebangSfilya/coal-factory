package equipment

import (
	"errors"
	"log"
	"strings"
)

const (
	PickCost    = 5
	VentCost    = 15_000
	TrolleyCost = 50_000
)

const (
	pick    = "pick"
	vent    = "vent"
	trolley = "trolley"
)

type Equipments struct {
	Pick    Pick
	Vent    Vent
	Trolley Trolley
	AllBuy  bool
}

func NewEquipmet() *Equipments {

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
	case pick:
		e.Pick.Buy()
		return pick, nil
	case vent:
		e.Vent.Buy()
		return vent, nil
	case trolley:
		e.Trolley.Buy()
		return trolley, nil
	default:
		log.Print("not buy")
		return "", errors.New("")
	}
}

func (e *Equipments) AllBuyed() bool {
	if e.Pick.IsBuyed && e.Trolley.IsBuyed && e.Vent.IsBuyed {
		e.AllBuy = true
		return e.AllBuy
	}
	return false
}

type Pick struct {
	IsBuyed bool
	Cost    int
}

func (e *Pick) Buy() {
	e.IsBuyed = true
}

type Vent struct {
	IsBuyed bool
	Cost    int
}

func (e *Vent) Buy() {
	e.IsBuyed = true
}

type Trolley struct {
	IsBuyed bool
	Cost    int
}

func (e *Trolley) Buy() {
	e.IsBuyed = true
}
