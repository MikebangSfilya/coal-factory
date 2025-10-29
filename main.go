package main

import (
	"coalFactory/api/handlers"
	"coalFactory/api/server"
	"coalFactory/config"
	"coalFactory/equipment"
	"coalFactory/factory"
	"coalFactory/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// the application is listening for the SIGTERM signal to exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	equipment.Init(cfg) //для тестирования, в будущем скорее всего уберу

	equip := equipment.NewEquipment()

	company := factory.Start(equip)

	service := service.New(company)

	handl := handlers.New(service)

	server := server.New(handl)

	go func() {
		if err := server.Start(); err != nil {
			log.Print("Server error:", err)
		}
	}()

	<-ctx.Done()

	company.CompanyStop()

	println("Shutting down gracefully...")
}
