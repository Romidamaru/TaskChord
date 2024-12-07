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
		api.GET("/task", func(c *gin.Context) {
			// Extract query parameters
			guildID := c.Query("guild_id")
			userID := c.Query("user_id")
			taskID := c.Query("task_id")

			// Call the controller method
			tasks, err := r.taskModel.Controller.GetTasksByUserID(guildID, userID, taskID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Respond with tasks
			c.JSON(http.StatusOK, tasks)
		})
	}
}
