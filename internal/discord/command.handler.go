package discord

import (
	"TaskChord/internal/pkg/task/ctrl"
	"github.com/bwmarrin/discordgo"
	"log"
)

type CommandHandler struct {
	taskController ctrl.TaskController
}

// NewCommandHandler creates a new instance of CommandHandler
func NewCommandHandler(taskController ctrl.TaskController) *CommandHandler {
	return &CommandHandler{taskController: taskController}
}

// HandleCommand processes the commands issued by users
func (h *CommandHandler) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		// Add task to the database using the task service
		err := h.taskController.CreateTask(userID, title, description, priority)
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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Task successfully created!",
			},
		})
	}
}
