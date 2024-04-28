package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"time-tracker-backend/account"
	"time-tracker-backend/manager"
	"time-tracker-backend/x/xjwt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	TaskManager    *manager.TaskManager
	AccountManager *account.AccountManager
}

type MobileRoutes struct {
	AccountManager *account.AccountManager
}

func SetupTaskRoutes(r *gin.Engine, taskManager *manager.TaskManager, accountManager *account.AccountManager) {
	taskRoutes := TaskRoutes{
		TaskManager:    taskManager,
		AccountManager: accountManager,
	}

	api := r.Group("/api/v1")

	taskGroup := api.Group("/tasks", taskRoutes.Authentication)
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

	api.POST("/login", taskRoutes.Login)
	api.POST("/create-user", taskRoutes.CreateUser)
	api.GET("/magic-link", taskRoutes.MagicLink)
}

func SetupMobileRoutes(r *gin.Engine, accountManager *account.AccountManager) {
	mobileRoutes := MobileRoutes{
		AccountManager: accountManager,
	}

	api := r.Group("/api/v1")

	mobileGroup := api.Group("/mobile")
	mobileGroup.POST("/magic-link", mobileRoutes.MagicLinkMobile)
}

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

type MagicLinkResponse struct {
	MagicLink string `json:"magicLink"`
}

type CreateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Login godoc
// @Summary      Perform login
// @Description  Authenticate user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        UserRequest  		body       CreateUserRequest  	true   "User infos"
// @Success      200 {object} MagicLinkResponse
// @Failure      404
// @Router       /v1/create-user [post]
func (tr *TaskRoutes) CreateUser(c *gin.Context) {
	u := &CreateUserRequest{}
	err := c.BindJSON(u)
	if err != nil {
		c.JSON(404, err)
		return
	}

	magicLink, err := tr.AccountManager.CreateUser(u.Name, u.Email)
	if err != nil {
		c.JSON(502, err.Error())
		return
	}

	c.JSON(200, MagicLinkResponse{MagicLink: magicLink})
}

type LoginRequest struct {
	Email string `json:"email,omitempty"`
}

// Login godoc
// @Summary      Perform login
// @Description  Authenticate user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        email  		body       LoginRequest  	true   "Email"
// @Success      200 {object} MagicLinkResponse
// @Failure      404
// @Router       /v1/login [post]
func (tr *TaskRoutes) Login(c *gin.Context) {
	u := &LoginRequest{}
	err := c.BindJSON(u)
	if err != nil {
		c.JSON(404, "Bad Request")
		return
	}

	magicLink, err := tr.AccountManager.Login(u.Email)
	if err != nil {
		c.JSON(404, err.Error())
		return
	}

	c.JSON(200, MagicLinkResponse{MagicLink: magicLink})
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

type MagicToken struct {
	MagicToken string `json:"magicToken,omitempty" form:"magicToken"`
}

func (tr *TaskRoutes) MagicLink(c *gin.Context) {
	u := &MagicToken{}
	err := c.BindQuery(u)
	if err != nil {
		c.JSON(404, "Bad Request")
		return
	}

	token, refreshToken, err := tr.AccountManager.MagicLink(u.MagicToken)
	if err != nil {
		c.JSON(401, err.Error())
		return
	}
	c.SetCookie("token", token, 9999, "/", "localhost", false, true)
	c.SetCookie("refreshToken", refreshToken, 9999, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "http://localhost:9000/board")
}

// Login godoc
// @Summary      Perform login
// @Description  Authenticate user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        token  		body       MagicToken  	true   "token"
// @Success      204
// @Failure      404
// @Router       /v1/mobile/magic-link [post]
func (tr *MobileRoutes) MagicLinkMobile(c *gin.Context) {
	u := &MagicToken{}
	err := c.BindQuery(u)
	if err != nil {
		c.JSON(404, "Bad Request")
		return
	}

	token, refreshToken, err := tr.AccountManager.MagicLink(u.MagicToken)
	if err != nil {
		c.JSON(401, err.Error())
		return
	}
	c.SetCookie("token", token, 9999, "/", "localhost", false, true)
	c.SetCookie("refreshToken", refreshToken, 9999, "/", "localhost", false, true)

	// c.Redirect(http.StatusFound, "http://localhost:9000/board")
	c.JSON(http.StatusNoContent, nil)
}

func (tr *TaskRoutes) Authentication(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, "token is empty")
		return
	}

	err = xjwt.VerifyToken(token)
	if err != nil {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(401, "refreshToken is empty")
			return
		}

		err = xjwt.VerifyToken(refreshToken)
		if err != nil {
			c.JSON(401, "refreshToken is invalid")
			return
		}

		token, refreshToken, err := tr.AccountManager.RefreshToken(refreshToken)
		if err != nil {
			c.JSON(401, "generate refreshtoken error")
			return
		}

		c.SetCookie("token", token, 9999, "/", "localhost", false, true)
		c.SetCookie("refreshToken", refreshToken, 9999, "/", "localhost", false, true)

	}

	c.Next()
}
