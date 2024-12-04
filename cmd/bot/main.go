package main

import (
	"github.com/joho/godotenv"
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"taskchord/internal/discord"
	"taskchord/internal/pkg/task/ctrl"
	taskEnt "taskchord/internal/pkg/task/ent"
	"taskchord/internal/pkg/task/svc"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the bot token and database DSN from environment variables
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Initialize PostgresDB (you can swap this out with other DB types later)
	database, err := gossiper.NewDB(
		gossiper.PostgresDB,
		dsn,
		true,
		[]any{taskEnt.Task{}},
	)
	if err != nil {
		log.Fatalf("Failed to create database instance: %v", err)
	}

	taskService := svc.NewTaskService(database)
	taskController := ctrl.NewTaskController(taskService)

	// Create command handler
	commandHandler := discord.NewCommandHandler(*taskController)

	// Create and start the bot
	bot, err := discord.NewBot(token, commandHandler)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	// Wait for termination signal to gracefully shut down the bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Block until a termination signal is received

	log.Println("Shutting down the bot...")
	bot.Stop()
}
