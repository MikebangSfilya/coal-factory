package config

import (
	"os"
	"strconv"
)

type Configuration struct {
	PickCost    int
	VentCost    int
	TrolleyCost int
}

func Load() *Configuration {
	pickCost := os.Getenv("PICK_COST")
	ventCost := os.Getenv("VENT_COST")
	trolleyCost := os.Getenv("TROLLEY_COST")

	pickCostInt, _ := strconv.Atoi(pickCost)
	ventCostInt, _ := strconv.Atoi(ventCost)
	trolleyCostInt, _ := strconv.Atoi(trolleyCost)

	return &Configuration{
		PickCost:    pickCostInt,
		VentCost:    ventCostInt,
		TrolleyCost: trolleyCostInt,
	}

}
