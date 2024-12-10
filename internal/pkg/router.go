package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskchord/internal/pkg/auth"
	"taskchord/internal/pkg/task"
	"taskchord/internal/pkg/user"
)

type Router struct {
	taskModel   *task.Module
	userModel   *user.Module
	oauthModule *auth.Module
}

func NewRouter() *Router {
	return &Router{
		taskModel:   task.New(),
		userModel:   user.New(),
		oauthModule: auth.New(),
	}
}

func (r *Router) InitREST(router *gin.Engine) {
	api := router.Group("/api")
	{
		// OAuth: Start the Discord authentication flow
		api.POST("/auth/discord", func(c *gin.Context) {
			// Generate Discord auth URL
			authURL, err := r.oauthModule.AuthController.GetAuthURL()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Respond with the auth URL
			c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
		})

		// Callback for Discord OAuth
		// Callback for Discord OAuth
		api.GET("/auth/discord/callback", func(c *gin.Context) {
			code := c.DefaultQuery("code", "")
			if code == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
				return
			}

			// Handle the callback with the provided code
			userInfo, err := r.oauthModule.AuthController.HandleAuthCallback(c.Writer, c.Request, code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Check if the user exists in the database
			user, err := r.userModel.Controller.GetUserByID(userInfo["id"].(string))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
				return
			}

			// If the user doesn't exist, save the user in the database
			if user == nil {
				err := r.userModel.Controller.AddUser(userInfo["id"].(string), userInfo["username"].(string), userInfo["email"].(string), userInfo["avatar"].(string))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
					return
				}
			}

			// Redirect the user to the frontend after successful authentication
			c.Redirect(http.StatusFound, "http://localhost:3000")
		})

		// Get tasks
		api.GET("/task", func(c *gin.Context) {
			guildID := c.Query("guild_id")
			userID := c.Query("user_id")
			taskID := c.Query("task_id")

			tasks, err := r.taskModel.Controller.GetTasksByUserID(guildID, userID, taskID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if len(tasks) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"message": "No tasks found for the provided criteria.",
					"tasks":   tasks, // Empty list to maintain consistency
				})
				return
			}

			c.JSON(http.StatusOK, tasks)
		})

		// Create task
		api.POST("/task", func(c *gin.Context) {
			var request struct {
				GuildID     string `json:"guild_id" binding:"required"`
				UserID      string `json:"user_id" binding:"required"`
				Title       string `json:"title" binding:"required"`
				Description string `json:"description"`
				Priority    string `json:"priority"`
				ExecutorID  string `json:"executor_id"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			taskID, err := r.taskModel.Controller.CreateTask(request.GuildID, request.UserID, request.Title, request.Description, request.Priority, request.ExecutorID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"task_id": taskID})
		})

		// Update task
		api.PUT("/task/:id", func(c *gin.Context) {
			taskID := c.Param("id")
			var request struct {
				GuildID     string `json:"guild_id" binding:"required"`
				UserID      string `json:"user_id" binding:"required"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Priority    string `json:"priority"`
				ExecutorID  string `json:"executor_id"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			updatedTaskID, err := r.taskModel.Controller.UpdateTask(request.GuildID, request.UserID, request.Title, request.Description, request.Priority, request.ExecutorID, taskID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"updated_task_id": updatedTaskID})
		})

		// Delete task
		api.DELETE("/task/:id", func(c *gin.Context) {
			taskID := c.Param("id")
			guildID := c.Query("guild_id")
			userID := c.Query("user_id")

			deletedTaskID, err := r.taskModel.Controller.DeleteTask(guildID, userID, taskID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"deleted_task_id": deletedTaskID})
		})
	}
}
