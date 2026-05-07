package main

import (
	app_worker "cloud/internal/app/worker"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Printf(".env parsing failed: %v", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
	)
	defer cancel()

	app, err := app_worker.New(ctx)
	if err != nil {
		fmt.Println("Failed to initialize worker app:", err)
		os.Exit(1)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		fmt.Printf("Worker app startup failed: %v", err)
		os.Exit(1)
	}
}
