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
	"time"
)

func main() {

	// the application is listening for the SIGTERM signal to exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	equipment.Init(cfg)

	equip := equipment.NewEquipment()

	company := factory.Start(ctx, equip)

	service := service.New(company)

	handl := handlers.New(service)

	server := server.New(handl)

	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			log.Print("Server error:", err)
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("recived shotdown signal")
	case err := <-errChan:
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	company.CompanyStop()

	log.Print("Shutdown end")

}
