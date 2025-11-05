package config

import (
	"os"
	"strconv"
)

type Configurate struct {
	PickCost    int
	VentCost    int
	TrolleyCost int
}

func Load() *Configurate {
	pickCost := os.Getenv("PICK_COST")
	ventCost := os.Getenv("VENT_COST")
	trolleyCost := os.Getenv("TROLLEY_COST")
	pickCostInt, _ := strconv.Atoi(pickCost)
	ventCostInt, _ := strconv.Atoi(ventCost)
	trolleyCostInt, _ := strconv.Atoi(trolleyCost)

	return &Configurate{
		PickCost:    pickCostInt,
		VentCost:    ventCostInt,
		TrolleyCost: trolleyCostInt,
	}

}
