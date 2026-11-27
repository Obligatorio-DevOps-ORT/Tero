package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

// requestLogger registra cada request HTTP como una línea JSON,
// usando el logger por defecto de slog (configurado en main.go).
// No registra /health para no llenar los logs con los health checks.
func requestLogger() gin.HandlerFunc {
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

// gormSlog es un "adaptador": GORM sabe escribir en cualquier cosa que tenga
// un método Printf, así que le damos uno que manda el mensaje a slog (JSON).
type gormSlog struct{}

func (gormSlog) Printf(format string, args ...interface{}) {
	slog.Warn("gorm", "detalle", fmt.Sprintf(format, args...))
}

// newGormLogger configura qué cosas registra GORM:
// solo errores reales y consultas lentas (más de 200 ms),
// sin el "record not found".
func newGormLogger() logger.Interface {
	return logger.New(gormSlog{}, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
}
