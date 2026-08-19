package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/config"
	"github.com/jj.jobo/FGC/internal/container"
	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/middleware"
	"github.com/jj.jobo/FGC/internal/routes"
)

func main() {

	// Load application configuration.
	config.LoadConfig()

	// Connect to PostgreSQL.
	database.ConnectDatabase()

	// Create Fiber application.
	app := fiber.New(fiber.Config{
		AppName:      config.App.AppName,
		ErrorHandler: middleware.ErrorHandler,
	})

	// Global middleware.
	app.Use(middleware.Recovery())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Dependency container.
	c := container.BuildContainer()

	// Routes.
	routes.Setup(app, c)

	// Start server.
	go func() {

		address := ":" + config.App.AppPort

		log.Printf(
			"%s running on %s",
			config.App.AppName,
			address,
		)

		if err := app.Listen(address); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	// Wait for shutdown signal.
	waitForShutdown(app)
}

func waitForShutdown(app *fiber.App) {

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-signalChannel

	log.Println("Shutdown signal received")

	shutdownTimeout := 10 * time.Second

	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		log.Printf("Fiber shutdown error: %v", err)
	}

	database.CloseDatabase()

	log.Println("Application shutdown complete")
}
