package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"taskchord/internal/pkg/task/ctrl"
)

// TODO
// 	add reminder and deadlines for tasks (optional), estimated time for task (optional)
// 	author and executor. Executor and Author can see same task if they bind to it

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
	case "delete":
		h.handleDeleteCommand(s, i)
	case "update":
		h.handleUpdateCommand(s, i)
	}
}

// HandleCreateCommand processes the commands issued by users
func (h *CommandHandler) handleCreateCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Extract required options
	title := options[0].StringValue()
	description := options[1].StringValue()

	// Set default values for optional options
	priority := "Medium"
	executorID := i.Interaction.Member.User.ID // Default to the creator

	// Process optional options dynamically
	for _, opt := range options[2:] { // Start processing optional arguments
		switch opt.Type {
		case discordgo.ApplicationCommandOptionString:
			priority = opt.StringValue() // Handle priority
		case discordgo.ApplicationCommandOptionUser:
			executorID = opt.UserValue(nil).ID // Handle executor
		}
	}

	// Create task
	userID := i.Member.User.ID
	guildID := i.GuildID

	taskIdInGuild, err := h.taskController.CreateTask(guildID, userID, title, description, priority, executorID)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create task. Please try again later.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Respond with success
	taskIDStr := strconv.FormatUint(uint64(taskIdInGuild), 10)
	embed := &discordgo.MessageEmbed{
		Color:       0x00FF00,
		Description: fmt.Sprintf("Task **#%s %s** successfully created!", taskIDStr, title),
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *CommandHandler) handleUpdateCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleUpdateCommand: %v", r)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "An unexpected error occurred while processing your request. Please try again.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}()

	userID := i.Interaction.Member.User.ID
	guildID := i.GuildID
	var id, title, description, priority string

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to update task. Please provide at least one field to update.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Extract the task ID (required)
	id = options[0].StringValue()

	// Set default values for optional fields
	priority = "Medium" // Default priority
	title = ""          // Default empty title
	description = ""    // Default empty description

	// Process optional fields dynamically
	for _, opt := range options[1:] { // Skip the first option (task ID)
		switch opt.Type {
		case discordgo.ApplicationCommandOptionString:
			// Handle title, description, and priority
			switch opt.Name {
			case "title":
				title = opt.StringValue()
			case "description":
				description = opt.StringValue()
			case "priority":
				priority = opt.StringValue()
			}
		}
	}

	// Validate and assign the title, description, and priority if they are provided
	if title == "" && description == "" && priority == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to update task. Please provide at least one field (title, description, or priority).",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Call the controller to update the task with the dynamically provided fields
	taskIdInGuild, err := h.taskController.UpdateTask(guildID, userID, title, description, priority, id)
	if err != nil {
		// Handle error based on the specific message
		var responseMessage string
		switch err.Error() {
		case "task ID is required":
			responseMessage = "Failed to update task. Task ID is required."
		case "invalid priority value":
			responseMessage = "Failed to update task. Priority value must be 'High', 'Medium', or 'Low'."
		default:
			responseMessage = "Failed to update task. Please try again later."
		}

		log.Printf("Error updating task: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: responseMessage,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Respond with the updated task info
	taskIDStr := strconv.FormatUint(uint64(taskIdInGuild), 10)
	embed := &discordgo.MessageEmbed{
		Color:       0x00FF00, // Green color
		Description: fmt.Sprintf("Task **#%s %s** successfully updated!", taskIDStr, title),
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

func GetNicknameFromID(userID string, s *discordgo.Session, guildID string) string {
	if userID == "" {
		return "Unknown User"
	}
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		if discordErr, ok := err.(*discordgo.RESTError); ok && discordErr.Response.StatusCode == 404 {
			log.Printf("User %s not found in guild %s.", userID, guildID)
			return "Unknown User"
		}
		log.Printf("Error fetching member details for userID %s: %v", userID, err)
		return "Unknown User"
	}
	if member.Nick != "" {
		return member.Nick
	}
	return member.User.Username // Fall back to username if nickname is not set
}

var nicknameCache = make(map[string]string)

func GetNicknameFromIDWithCache(userID string, s *discordgo.Session, guildID string) string {
	if nickname, found := nicknameCache[userID]; found {
		return nickname
	}

	nickname := GetNicknameFromID(userID, s, guildID)
	nicknameCache[userID] = nickname
	return nickname
}

func (h *CommandHandler) handleShowCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Interaction.Member.User.ID
	guildID := i.GuildID
	var id string

	options := i.ApplicationCommandData().Options
	if len(options) > 0 { // Check if the id option is provided
		id = options[0].StringValue() // Guaranteed to be "High", "Medium", or "Low" from the select menu
	}

	// Retrieve tasks from the database
	tasks, err := h.taskController.GetTasksByUserID(guildID, userID, id)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch tasks. Please try again later.",
				Flags:   discordgo.MessageFlagsEphemeral,
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
			taskIDStr := strconv.FormatUint(uint64(task.TaskIdInGuild), 10)

			// Use cached nickname retrieval
			authorNickname := GetNicknameFromIDWithCache(task.UserID, s, guildID)
			executorNickname := GetNicknameFromIDWithCache(task.ExecutorID, s, guildID)

			description := fmt.Sprintf(
				"Author: <@%s> (%s)\nExecutor: <@%s> (%s)\nPriority: %s\n**Description:**\n%s",
				task.UserID, authorNickname,
				task.ExecutorID, executorNickname,
				string(task.Priority), task.Description,
			)

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "**#" + taskIDStr + " " + task.Title + "**",
				Value:  description,
				Inline: false,
			})

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
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *CommandHandler) handleDeleteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Interaction.Member.User.ID
	guildID := i.GuildID
	var id string

	// Retrieve the task ID from the options
	options := i.ApplicationCommandData().Options
	if len(options) > 0 {
		id = options[0].StringValue()
	}

	// Call the controller to delete the task
	taskIdInGuild, err := h.taskController.DeleteTask(guildID, userID, id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to delete task. Please try again later.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Respond to the user with confirmation
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Task #%s successfully deleted!", taskIdInGuild),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
