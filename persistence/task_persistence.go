package persistence

import (
	"time-tracker-backend/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Cria uma task no banco de dados
func (tr *TaskRepository) CreateTask(task *models.Task) error {
	result := tr.db.Create(task)
	return result.Error
}

func (tr *TaskRepository) CreateTimeRecord(timeRecord *models.RegisteredTime) error {
	result := tr.db.Create(timeRecord)
	return result.Error
}

// Lista todas as tasks no banco de dados (sem registros excluídos)
func (tr *TaskRepository) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	result := tr.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (tr *TaskRepository) DeleteTask(userID, taskID string) error {
	var task models.Task
	result := tr.db.First(&task, "id = ? AND user_id = ?", taskID, userID)
	if result.Error != nil {
		return result.Error
	}

	result = tr.db.Delete(&task)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	result = tr.db.Model(&models.RegisteredTime{}).Where("task_id = ?", taskID).Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return result.Error
	}

	return nil
}
