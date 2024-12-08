package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskchord/internal/pkg/task"
)

type Router struct {
	taskModel *task.Module
}

func NewRouter() *Router {
	return &Router{
		taskModel: task.New(),
	}
}

func (r *Router) InitREST(router *gin.Engine) {
	api := router.Group("/api")
	{
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
