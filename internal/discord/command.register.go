package discord

import "github.com/bwmarrin/discordgo"

func RegisterCommands(s *discordgo.Session) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create",
			Description: "Create a new task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Title of the task",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Description of the task",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "priority",
					Description: "Priority of the task",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "High",
							Value: "High",
						},
						{
							Name:  "Medium",
							Value: "Medium",
						},
						{
							Name:  "Low",
							Value: "Low",
						},
					},
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
