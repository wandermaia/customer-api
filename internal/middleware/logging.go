package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger é um middleware que registra informações sobre as requisições HTTP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tempo de início
		startTime := time.Now()

		// Processa a requisição
		c.Next()

		// Tempo de término
		endTime := time.Now()

		// Tempo de execução
		latency := endTime.Sub(startTime)

		// Detalhes da requisição
		path := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// Log
		log.Printf("[%s] %s %s %d %s", method, path, latency, statusCode, c.ClientIP())
	}
}
