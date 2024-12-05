package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"taskchord/internal/pkg/task/ctrl"
)

// TODO
// delete task by user id
// show only one task by id
// bind task to guild id (discord)

type CommandHandler struct {
	taskController ctrl.TaskController
}

// NewCommandHandler creates a new instance of CommandHandler
func NewCommandHandler(taskController ctrl.TaskController) *CommandHandler {
	return &CommandHandler{taskController: taskController}
}

// HandleCommand processes the commands issued by users
func (h *CommandHandler) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "create":
		h.handleCreateCommand(s, i)
	case "show":
		h.handleShowCommand(s, i)
	}
}

// HandleCreateCommand processes the commands issued by users
func (h *CommandHandler) handleCreateCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Ensure we're handling the `/create` command
	if i.ApplicationCommandData().Name == "create" {
		options := i.ApplicationCommandData().Options
		title := options[0].StringValue()
		description := options[1].StringValue()

		priority := "Medium"  // Default value
		if len(options) > 2 { // Check if the priority option is provided
			priority = options[2].StringValue() // Guaranteed to be "High", "Medium", or "Low" from the select menu
		}

		userID := i.Interaction.Member.User.ID
		guildID := i.GuildID

		// Add task to the database using the task service
		taskIdInGuild, err := h.taskController.CreateTask(guildID, userID, title, description, priority)
		if err != nil {
			log.Printf("Error creating task: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to create task. Please try again later.",
				},
			})
			return
		}

		// Respond to the user
		taskIDStr := strconv.FormatUint(uint64(taskIdInGuild), 10)

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Color:       0x00FF00, // Green color
			Description: "Task **#" + taskIDStr + " " + title + "** successfully created!",
		}

		// Respond to the user with the embed
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}
}

func (h *CommandHandler) handleShowCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Interaction.Member.User.ID
	guildID := i.GuildID

	// Retrieve tasks from the database
	tasks, err := h.taskController.GetTasksByUserID(guildID, userID)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch tasks. Please try again later.",
			},
		})
		return
	}

	// Create embed response
	embed := &discordgo.MessageEmbed{
		Title:  "Your Tasks:",
		Color:  0x00FF00, // Green color
		Fields: []*discordgo.MessageEmbedField{},
	}

	if len(tasks) == 0 {
		embed.Description = "You have no tasks!"
	} else {
		for i, task := range tasks {
			// Convert task.ID (uint) to string
			taskIDStr := strconv.FormatUint(uint64(task.TaskIdInGuild), 10)

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "**#" + taskIDStr + " " + task.Title + "**", // Include task ID in the field name
				Value:  "Priority: " + string(task.Priority) + "\n**Description:**\n" + task.Description,
				Inline: false,
			})

			// Add a separator for all tasks except the last one
			if i < len(tasks)-1 {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "\u200B",
					Value:  "──────────────",
					Inline: false,
				})
			}
		}
	}

	// Respond with the embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
