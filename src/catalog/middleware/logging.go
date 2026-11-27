package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger registra cada request HTTP como una línea JSON,
// usando el logger por defecto de slog (configurado en main.go).
// No registra /health para no llenar los logs con los health checks.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/health" {
			c.Next()
			return
		}

		inicio := time.Now()
		c.Next() // acá Gin ejecuta el handler de la ruta

		status := c.Writer.Status()
		nivel := slog.LevelInfo
		if status >= 500 {
			nivel = slog.LevelError
		} else if status >= 400 {
			nivel = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), nivel, "request",
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration_ms", float64(time.Since(inicio).Microseconds())/1000,
		)
	}
}
