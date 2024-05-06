package routes

import (
	"log"
	"time-tracker-backend/account"
	"time-tracker-backend/config"
	manager "time-tracker-backend/controllers"
	"time-tracker-backend/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	database.ConnectDataBase()

	db, err := database.ConnectDataBase()
	if err != nil {
		log.Fatal("Erro ao conectar " + err.Error())
	}

	// Criando o gerenciador de tarefas
	taskManager := manager.NewTaskManager(db)

	p := account.NewPersistence(db)
	accountManager := account.NewManager(p)

	r := gin.Default()

	// Middleware para habilitar o CORS. Remover no futuro
	cfg := cors.DefaultConfig()
	// config.AllowAllOrigins = true // no set *
	cfg.AllowCredentials = true
	cfg.AllowOrigins = []string{"http://localhost:9000"}
	r.Use(cors.New(cfg))
	SetupTaskRoutes(r, taskManager, accountManager)
	// SetupMobileRoutes(r, accountManager)
	r.Run(config.HttpPort())
}
