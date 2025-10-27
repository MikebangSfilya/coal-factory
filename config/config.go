package config

type Configurate struct {
	PickCost    int
	VentCost    int
	TrolleyCost int
}

func Load() *Configurate {
	return &Configurate{
		PickCost:    3_000,
		VentCost:    15_000,
		TrolleyCost: 50_000,
	}

}
