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

	"github.com/joho/godotenv"
)

// @title           Coal Factory API
// @version         1.0
// @description     API coal-factory

// @host      localhost:8080
// @BasePath  /
func main() {

	// the application is listening for the SIGTERM signal to exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Printf(".env not found: %v", err)
	}

	cfg := config.Load()

	equipment.Load(cfg)

	equip := equipment.NewEquipment()

	company := factory.Start(ctx, equip)

	gameService := service.New(company)

	handle := handlers.New(gameService)

	srv := server.New(":8080", handle)

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil {
			log.Print("Server error:", err)
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP srv shutdown error: %v", err)
	}
	company.CompanyStop()

	log.Print("Shutdown end")

}
