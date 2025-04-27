package main

import (
	"log"

	"github.com/wandermaia/customer-api/docs"
	"github.com/wandermaia/customer-api/internal/config"
	"github.com/wandermaia/customer-api/internal/domain/repository"
	"github.com/wandermaia/customer-api/internal/domain/service"
	"github.com/wandermaia/customer-api/internal/handler"
	"github.com/wandermaia/customer-api/internal/middleware"
	"github.com/wandermaia/customer-api/pkg/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cliente API
// @version 1.0
// @description API RESTful para gestão de clientes
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	// Configura informações do Swagger
	docs.SwaggerInfo.Title = "Cliente API"
	docs.SwaggerInfo.Description = "API RESTful para gestão de clientes em Go usando Gin Framework"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Carrega as configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Falha ao carregar as configurações: %v", err)
	}

	// Configura o modo do Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializa o banco de dados
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Inicializa o repositório
	customerRepo := repository.NewPostgresCustomerRepository(db)

	// Inicializa o serviço
	customerService := service.NewCustomerService(customerRepo)

	// Inicializa o handler
	customerHandler := handler.NewCustomerHandler(customerService)

	// Configura o router
	router := gin.Default()

	// Adiciona middleware
	router.Use(middleware.Logger())

	// Registra as rotas
	customerHandler.RegisterRoutes(router)

	// Adiciona rota de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Adiciona o endpoint do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicia o servidor
	log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
	log.Printf("Documentação Swagger disponível em http://localhost:%s/swagger/index.html", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
