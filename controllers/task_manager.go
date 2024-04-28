package manager

import (
	"time-tracker-backend/models"
	"time-tracker-backend/persistence"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type TaskManager struct {
	db             *gorm.DB
	taskRepository *persistence.TaskRepository
}

func NewTaskManager(db *gorm.DB) *TaskManager {
	return &TaskManager{
		db:             db,
		taskRepository: persistence.NewTaskRepository(db),
	}
}

func (tm *TaskManager) CreateTask(userID, title, tag string) (*models.Task, error) {
	task := &models.Task{
		ID:             uuid.NewV4().String(),
		UserID:         userID,
		Title:          title,
		Tag:            tag,
		RegisteredTime: []models.RegisteredTime{},
	}

	result := tm.db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (tm *TaskManager) ListTasks(c *gin.Context) ([]models.Task, error) {
	var tasks []models.Task

	result := tm.db.Preload("RegisteredTime", "end_time IS NOT NULL").
		Where("done = false AND deleted_at IS NULL").
		Order("updated_at desc").
		Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (tm *TaskManager) ListDoneTasks(c *gin.Context) ([]models.Task, error) {
	var doneTasks []models.Task

	result := tm.db.Where("done = ?", true).Find(&doneTasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return doneTasks, nil
}

func (tm *TaskManager) DeleteTask(userID, taskID string) error {
	return tm.taskRepository.DeleteTask(userID, taskID)
}

func (tm *TaskManager) CompleteTask(userID, taskID string) (*models.Task, error) {
	var task models.Task

	result := tm.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	task.Done = true

	result = tm.db.Save(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

func (tm *TaskManager) AddTagToTask(userID, taskID, tag string) error {
	var task models.Task
	result := tm.db.First(&task, "id = ? AND user_id = ?", taskID, userID)
	if result.Error != nil {
		return result.Error
	}

	task.Tag = tag

	result = tm.db.Save(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (tm *TaskManager) RemoveTagFromTask(userID, taskID string) error {
	var task models.Task
	result := tm.db.First(&task, "id = ? AND user_id = ?", taskID, userID)
	if result.Error != nil {
		return result.Error
	}

	task.Tag = ""

	result = tm.db.Save(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (tm *TaskManager) UpdateTaskTitle(userID, taskID, newTitle string) (*models.Task, error) {
	var task models.Task
	result := tm.db.First(&task, "id = ? AND user_id = ?", taskID, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	task.Title = newTitle

	result = tm.db.Save(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}
