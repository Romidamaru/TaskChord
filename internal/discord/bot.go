package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	Token          string
	CommandHandler *CommandHandler
	Session        *discordgo.Session
}

// NewBot creates a new bot instance
func NewBot(token string, handler *CommandHandler) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:          token,
		CommandHandler: handler,
		Session:        session,
	}, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	// Add the command handler
	b.Session.AddHandler(b.CommandHandler.HandleCommand)

	// Open the bot session
	err := b.Session.Open()
	if err != nil {
		return err
	}

	//// Delete old commands
	//log.Println("Bot session opened. Cleaning up old commands...")
	//err = b.DeleteCommands()
	//if err != nil {
	//	log.Printf("Failed to delete old commands: %v", err)
	//} else {
	//	log.Println("Old commands deleted successfully.")
	//}
	//
	//// Register commands
	//log.Println("Bot session opened. Registering commands...")
	//err = RegisterCommands(b.Session)
	//if err != nil {
	//	log.Printf("Failed to register commands: %v", err)
	//} else {
	//	log.Println("Commands registered successfully.")
	//}

	log.Println("\033[32m" + "Bot is running..." + "\033[0m")

	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.Session.Close()
	log.Println("Bot stopped.")
}
