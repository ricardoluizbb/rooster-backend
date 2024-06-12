package routes

import (
	"context"
	"log"
	"time-tracker-backend/config"
	"time-tracker-backend/controllers"
	manager "time-tracker-backend/controllers"
	"time-tracker-backend/database"
	"time-tracker-backend/models"

	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func SetupRoutes() {
	database.ConnectDataBase()

	db, err := database.ConnectDataBase()
	if err != nil {
		log.Fatal("Erro ao conectar " + err.Error())
	}

	// Crie uma instância do app Firebase
	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	// Crie uma instância do cliente de autenticação Firebase
	fireAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	// Crie uma instância do controlador de autenticação usando o construtor
	authService := &models.AuthService{
		DB:       db,
		FireAuth: fireAuth,
	}
	authController := controllers.NewAuthController(authService)

	// Criando o gerenciador de tarefas
	taskManager := manager.NewTaskManager(db)

	r := gin.Default()

	// Middleware para habilitar o CORS. Remover no futuro
	cfg := cors.DefaultConfig()
	// config.AllowAllOrigins = true // no set *
	cfg.AllowCredentials = true
	cfg.AllowOrigins = []string{"http://localhost:9000"}
	r.Use(cors.New(cfg))
	SetupTaskRoutes(r, taskManager, authController)
	r.Run(config.HttpPort())
}
