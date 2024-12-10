package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"taskchord/internal/pkg/task"

	"github.com/gin-gonic/gin"
	"taskchord/internal/core/config"
	"taskchord/internal/discord"
	"taskchord/internal/pkg"
)

func main() {
	token := config.Inst().DiscordToken

	taskModule := task.New()

	// Create Discord bot
	commandHandler := discord.NewCommandHandler(*taskModule.Controller)
	bot, err := discord.NewBot(token, commandHandler)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Create REST server
	router := gin.Default()
	appRouter := pkg.NewRouter()
	appRouter.InitREST(router)

	// Use WaitGroup to manage parallel execution
	var wg sync.WaitGroup

	// Start REST API server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting REST API server on :8080...")
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	// Start Discord bot
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting Discord bot...")
		if err := bot.Start(); err != nil {
			log.Fatalf("Failed to start Discord bot: %v", err)
		}
	}()

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Block until a termination signal is received

	log.Println("Shutting down the bot...")
	bot.Stop()

	log.Println("Shutting down the REST API server...")

	// Wait for all goroutines to complete
	wg.Wait()

	log.Println("Application stopped gracefully.")
}
