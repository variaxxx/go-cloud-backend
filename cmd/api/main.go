package main

import (
	app_api "cloud/internal/app/api"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
	)
	defer cancel()

	app, err := app_api.New()
	if err != nil {
		fmt.Println("Failed to initialize app:", err)
		os.Exit(1)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		fmt.Printf("HTTP server startup failed: %v", err)
		os.Exit(1)
	}
}
