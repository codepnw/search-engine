package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codepnw/search-engine/internal/api"
	"github.com/codepnw/search-engine/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

const envFile = "test.env"

func main() {
	if err := godotenv.Load(envFile); err != nil {
		panic("cannot find .env file")
	}

	var port string
	if port = os.Getenv("APP_PORT"); port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())
	db.InitDB()  // Database
	api.NewRoutes(app) // API Routes

	// Start and Shutdown Server
	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block the main thread untill interupted
	app.Shutdown()
	fmt.Println("shutting down server...")
}
