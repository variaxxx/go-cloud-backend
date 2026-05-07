package main

import (
	app_load "cloud/internal/app/load"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Printf(".env parsing failed: %v", err)
		os.Exit(1)
	}

	app, err := app_load.New()
	if err != nil {
		fmt.Println("Failed to initialize load app:", err)
		os.Exit(1)
	}
	defer app.Close()

	if err := app.Run(context.Background()); err != nil {
		fmt.Printf("Load app failed: %v", err)
		os.Exit(1)
	}
}
