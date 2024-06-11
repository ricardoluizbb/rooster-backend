package routes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func (tr *TaskRoutes) ListTasks(c *gin.Context) {
	tasks, err := tr.TaskManager.ListTasks(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tasks)
}

func (tr *TaskRoutes) DeleteTask(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	err := tr.TaskManager.DeleteTask(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Tarefa excluída com sucesso"})
}

func (tr *TaskRoutes) ListDoneTasks(c *gin.Context) {
	doneTasks, err := tr.TaskManager.ListDoneTasks(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, doneTasks)
}

func (tr *TaskRoutes) CreateTask(c *gin.Context) {
	var request struct {
		Title string `json:"title"`
		Tag   string `json:"tag"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := "userID"

	task, err := tr.TaskManager.CreateTask(userID, request.Title, request.Tag)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, task)
}

func (tr *TaskRoutes) StartTask(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	err := tr.TaskManager.StartTask(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Timer iniciado"})
}

func (tr *TaskRoutes) PauseTask(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	timeRecord, err := tr.TaskManager.PauseTask(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"time_record": timeRecord})
}

func (tr *TaskRoutes) AddTagToTask(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	var request struct {
		Tag string `json:"tag"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := tr.TaskManager.AddTagToTask(userID, taskID, request.Tag)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Tag adicionada com sucesso"})
}

func (tr *TaskRoutes) UpdateTaskTitle(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	var request struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedTask, err := tr.TaskManager.UpdateTaskTitle(userID, taskID, request.Title)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, updatedTask)
}

func (tr *TaskRoutes) GetTaskTotalTime(c *gin.Context) {
	taskID := c.Param("taskID")

	totalTime, err := tr.TaskManager.GetTaskTotalTime(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, totalTime)
}

func (tr *TaskRoutes) UpdateRegisteredTime(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")
	registeredTimeID := c.Param("registeredTimeID")

	var request struct {
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedRegisteredTime, err := tr.TaskManager.EditRegisteredTime(userID, taskID, registeredTimeID, request.StartTime, request.EndTime)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, updatedRegisteredTime)
}

func (tr *TaskRoutes) CompleteTask(c *gin.Context) {
	userID := "userID"
	taskID := c.Param("taskID")

	completedTask, err := tr.TaskManager.CompleteTask(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, completedTask)
}

func (tr *TaskRoutes) ExportTaskTotalTimesToExcel(c *gin.Context) {
	tasks, err := tr.TaskManager.ListTasks(nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	doneTasks, err := tr.TaskManager.ListDoneTasks(nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	allTasks := append(tasks, doneTasks...)

	file := excelize.NewFile()
	sheetName := "Relatório"
	index := file.NewSheet(sheetName)

	file.DeleteSheet("Sheet1")

	headers := map[string]string{
		"A1": "Titulo",
		"B1": "Início",
		"C1": "Fim",
		"D1": "Tempo total",
		"E1": "Concluida",
	}

	style, err := file.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for cell, header := range headers {
		file.SetCellValue(sheetName, cell, header)
		file.SetCellStyle(sheetName, cell, cell, style)
	}

	row := 2
	for _, task := range allTasks {
		totalTime, err := tr.TaskManager.GetTaskTotalTime(nil, task.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		file.SetCellValue(sheetName, "A"+strconv.Itoa(row), task.Title)

		if len(totalTime.RegisteredTimes) > 0 {
			startTimeFormatted := totalTime.RegisteredTimes[0].StartTime.Format("02/01/2006 15:04:05")
			file.SetCellValue(sheetName, "B"+strconv.Itoa(row), startTimeFormatted)
		}

		if len(totalTime.RegisteredTimes) > 0 {
			lastIndex := len(totalTime.RegisteredTimes) - 1
			endTimeFormatted := totalTime.RegisteredTimes[lastIndex].EndTime.Format("02/01/2006 15:04:05")
			file.SetCellValue(sheetName, "C"+strconv.Itoa(row), endTimeFormatted)
		}

		totalMilliseconds := int(totalTime.TotalTime)
		totalSeconds := totalMilliseconds / 1000
		hours := totalSeconds / 3600
		minutes := (totalSeconds % 3600) / 60
		seconds := totalSeconds % 60
		formattedTotalTime := fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
		file.SetCellValue(sheetName, "D"+strconv.Itoa(row), formattedTotalTime)

		concluida := "Não"
		if task.Done {
			concluida = "Sim"
		}
		file.SetCellValue(sheetName, "E"+strconv.Itoa(row), concluida)

		row++
	}

	file.SetColWidth(sheetName, "A", "A", 50)
	file.SetColWidth(sheetName, "B", "E", 20)
	file.SetActiveSheet(index)

	filePath := "Relatorio.xlsx"
	if err := file.SaveAs(filePath); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}

func OptionsHandler(c *gin.Context) {
	c.JSON(200, gin.H{})
}
