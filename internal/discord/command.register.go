package discord

import "github.com/bwmarrin/discordgo"

func RegisterCommands(s *discordgo.Session) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create",
			Description: "Create a new task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "title",
					Description: "Title of the task",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "description",
					Description: "Description of the task",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "priority",
					Description: "Priority of the task (High, Medium, Low)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
	}

	// Register the commands
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd) // "" for global command; replace with Guild ID for guild-specific
		if err != nil {
			return err
		}
	}
	return nil
}
