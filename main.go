package main

import (
	"coalFactory/api/handlers"
	"coalFactory/api/server"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/service"
)

func main() {
	//1. Эквик
	//2. Компания
	//3. Сервис
	//4. Хендлеры
	//5. Сервер

	equip := equipment.NewEquipmet()

	company := factory.Start(equip)

	service := service.New(company, equip)

	handl := handlers.New(service)

	server := server.New(handl)

	server.Start()

}
