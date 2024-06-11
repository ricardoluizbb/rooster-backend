package controllers

import (
	"context"
	"errors"
	"time"
	"time-tracker-backend/models"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (tm *TaskManager) CreateRegisteredTime(userID string, task *models.Task, paused bool) error {
	registeredTime := &models.RegisteredTime{
		ID:        uuid.NewV4().String(),
		TaskID:    task.ID,
		Paused:    paused,
		UserID:    userID,
		StartTime: time.Now(),
	}

	result := tm.db.Create(registeredTime)
	return result.Error
}

func (tm *TaskManager) StartTask(userID, taskID string) error {
	var registeredTime models.RegisteredTime

	// Encontrar o registro mais recente usando Order By updatedAt
	result := tm.db.Where("task_id = ? AND user_id = ?", taskID, userID).
		Order("updated_at desc").First(&registeredTime)

	if result.Error != nil {
		// Se o registro não foi encontrado, crie um novo
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newID := uuid.NewV4().String()
			newRegisteredTime := models.RegisteredTime{
				ID:        newID,
				TaskID:    taskID,
				Paused:    false,
				UserID:    userID,
				StartTime: time.Now(),
				EndTime:   nil,
			}

			result = tm.db.Create(&newRegisteredTime)
			return result.Error
		}

		return result.Error
	}

	// Verificar se o atributo Paused está true
	if registeredTime.Paused {
		// Criar um novo registro
		newID := uuid.NewV4().String()
		newRegisteredTime := models.RegisteredTime{
			ID:        newID,
			TaskID:    taskID,
			Paused:    false,
			UserID:    userID,
			StartTime: time.Now(),
			EndTime:   nil,
		}

		result = tm.db.Create(&newRegisteredTime)
		return result.Error
	}
	return nil
}

func (tm *TaskManager) PauseTask(userID, taskID string) (*models.RegisteredTime, error) {
	var registeredTime models.RegisteredTime
	result := tm.db.Where("task_id = ? AND user_id = ? AND paused = ?", taskID, userID, false).First(&registeredTime)
	if result.Error != nil {
		return nil, result.Error
	}

	// Caso a tarefa já esteja pausada, não fazer nada
	if registeredTime.Paused {
		return nil, nil
	}

	currentTime := time.Now()
	registeredTime.EndTime = &currentTime
	registeredTime.Paused = true

	// Inicie uma transação
	tx := tm.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Atualize o registro existente na transação
	result = tx.Save(&registeredTime)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Commit da transação
	tx.Commit()

	// Chame GetTaskTotalTime para atualizar o tempo total da tarefa
	_, err := tm.GetTaskTotalTime(context.Background(), taskID)
	if err != nil {
		return nil, err
	}

	return &registeredTime, nil
}

func (tm *TaskManager) GetTaskTotalTime(ctx context.Context, taskID string) (*models.TaskTotalTime, error) {

	arr := []*models.RegisteredTime{}

	tx := tm.db.Model(&models.RegisteredTime{}).Where("task_id = ?", taskID).Find(&arr)
	if tx.Error != nil {
		return nil, tx.Error
	}

	sum := 0.0

	for _, rt := range arr {
		rt.SetTotalTime()
		sum += rt.TotalTime
	}

	return &models.TaskTotalTime{
		TaskID:          taskID,
		RegisteredTimes: arr,
		TotalTime:       sum,
	}, nil
}

func (tm *TaskManager) EditRegisteredTime(userID, taskID, registeredTimeID string, startTime, endTime time.Time) (*models.RegisteredTime, error) {
	var registeredTime models.RegisteredTime

	// Verificar se o registro pertence ao usuário e à tarefa especificados
	result := tm.db.Where("id = ? AND task_id = ? AND user_id = ?", registeredTimeID, taskID, userID).First(&registeredTime)
	if result.Error != nil {
		return nil, result.Error
	}

	// Atualizar StartTime e EndTime
	registeredTime.StartTime = startTime
	registeredTime.EndTime = &endTime

	// Iniciar uma transação
	tx := tm.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Salvar o registro atualizado na transação
	result = tx.Save(&registeredTime)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Commit da transação
	tx.Commit()

	// Chame GetTaskTotalTime para atualizar o tempo total da tarefa
	_, err := tm.GetTaskTotalTime(context.Background(), taskID)
	if err != nil {
		return nil, err
	}

	return &registeredTime, nil
}
