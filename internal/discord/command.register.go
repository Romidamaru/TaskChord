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
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "executor",
					Description: "Executor of the task",
					Required:    false,
				},
			},
		},
		{
			Name:        "show",
			Description: "Show all your tasks",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID of task",
					Required:    false,
				},
			},
		},
		{
			Name:        "update",
			Description: "Update a task by id",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID of task",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Title of the task",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Description of the task",
					Required:    false,
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
		{
			Name:        "delete",
			Description: "Delete your task by ID",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID of task",
					Required:    true,
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
