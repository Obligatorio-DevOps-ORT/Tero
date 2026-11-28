package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
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

// RecuperarPanico reemplaza a gin.Recovery(): si un handler entra en pánico,
// responde 500 y lo registra en JSON. gin.Recovery() lo escribe en texto,
// en varias líneas y con códigos de color.
// Va después de RequestLogger, para que el logger registre ese 500.
func RecuperarPanico() gin.HandlerFunc {
	// nil: que Gin no escriba nada por su cuenta; lo escribimos nosotros.
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err any) {
		slog.Error("pánico en un handler", "error", fmt.Sprint(err), "stack", string(debug.Stack()))
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
