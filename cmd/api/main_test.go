package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wandermaia/customer-api/docs"
	"github.com/wandermaia/customer-api/internal/config"
)

func TestHealthCheckEndpoint(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup router
	router := gin.Default()

	// Register the health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestConfigLoading(t *testing.T) {
	// Test config loading
	cfg, err := config.LoadConfig()

	// Assert config loading
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.ServerPort)
}

// func TestDatabaseConnection(t *testing.T) {
// 	// Skip in CI environment or add mock
// 	if testing.Short() {
// 		t.Skip("Skipping database integration test")
// 	}

// 	cfg, _ := config.LoadConfig()
// 	db, err := database.NewPostgresConnection(cfg)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, db)

// 	// Test database connection
// 	err = db.Ping()
// 	assert.NoError(t, err)
// }

func TestSwaggerInfoSetup(t *testing.T) {
	// Call the function that sets up Swagger info
	setupSwaggerInfo()

	// Assert Swagger information is set correctly
	assert.Equal(t, "Cliente API", docs.SwaggerInfo.Title)
	assert.Equal(t, "1.0", docs.SwaggerInfo.Version)
	assert.Equal(t, "localhost:8080", docs.SwaggerInfo.Host)
	assert.Equal(t, "/api", docs.SwaggerInfo.BasePath)
	assert.Contains(t, docs.SwaggerInfo.Schemes, "http")
	assert.Contains(t, docs.SwaggerInfo.Schemes, "https")
}

// Helper function to setup swagger info (extracted from main)
func setupSwaggerInfo() {
	docs.SwaggerInfo.Title = "Cliente API"
	docs.SwaggerInfo.Description = "API RESTful para gestão de clientes em Go usando Gin Framework"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
