package routes

import (
	"time-tracker-backend/controllers"
	manager "time-tracker-backend/controllers"

	"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	TaskManager    *manager.TaskManager
	AuthController *controllers.AuthController
}

func SetupTaskRoutes(r *gin.Engine, taskManager *manager.TaskManager, authController *controllers.AuthController) {
	taskRoutes := TaskRoutes{
		TaskManager:    taskManager,
		AuthController: authController,
	}

	api := r.Group("/api/v1")

	// Set up the authentication routes
	api.POST("/login", authController.Login)
	api.POST("/register", authController.Register)

	taskGroup := api.Group("/tasks")
	taskGroup.GET("", taskRoutes.ListTasks)
	taskGroup.GET("/done-tasks", taskRoutes.ListDoneTasks)
	taskGroup.POST("", taskRoutes.CreateTask)
	taskGroup.POST("/:taskID/start", taskRoutes.StartTask)
	taskGroup.POST("/:taskID/pause", taskRoutes.PauseTask)
	taskGroup.DELETE("/:taskID", taskRoutes.DeleteTask)
	taskGroup.PUT("/:taskID/complete", taskRoutes.CompleteTask)
	taskGroup.POST("/:taskID/tags", taskRoutes.AddTagToTask)
	taskGroup.PUT("/:taskID", taskRoutes.UpdateTaskTitle)
	taskGroup.GET("/:taskID/total-time", taskRoutes.GetTaskTotalTime)
	taskGroup.PUT("/:taskID/registered-times/:registeredTimeID", taskRoutes.UpdateRegisteredTime)
	taskGroup.GET("/export-to-excel", taskRoutes.ExportTaskTotalTimesToExcel)
	taskGroup.OPTIONS("", OptionsHandler)
}
