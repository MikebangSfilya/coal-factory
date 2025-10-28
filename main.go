package main

import (
	"coalFactory/api/handlers"
	"coalFactory/api/server"
	"coalFactory/config"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/service"
)

func main() {

	cfg := config.Load()

	equipment.Init(cfg) //для тестирования, в будущем скорее всего уберу

	equip := equipment.NewEquipment()

	company := factory.Start(equip)

	service := service.New(company)

	handl := handlers.New(service)

	server := server.New(handl)

	server.Start()

}
