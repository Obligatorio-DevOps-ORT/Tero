package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
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

// recuperarPanico reemplaza a gin.Recovery(): si un handler entra en pánico,
// responde 500 y lo registra en JSON. gin.Recovery() lo escribe en texto,
// en varias líneas y con códigos de color.
// Va después de requestLogger, para que el logger registre ese 500.
func recuperarPanico() gin.HandlerFunc {
	// nil: que Gin no escriba nada por su cuenta; lo escribimos nosotros.
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err any) {
		slog.Error("pánico en un handler", "error", fmt.Sprint(err), "stack", string(debug.Stack()))
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// gormLogger reemplaza al logger que trae GORM y escribe en slog (JSON),
// con el nivel que corresponde a cada caso:
//   - ERROR: un error de SQL o de conexión.
//   - WARN: una consulta que tardó más de 200 ms.
//
// No registra "record not found"
// ni los valores de las consultas: la consulta se escribe con $1, $2...
// porque los valores pueden traer datos de los clientes.
type gormLogger struct{}

func newGormLogger() logger.Interface {
	return gormLogger{}
}

// LogMode lo pide la interfaz de GORM; el nivel lo decidimos nosotros en Trace.
func (g gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return g
}

// Info: no registramos los mensajes informativos de GORM.
func (gormLogger) Info(context.Context, string, ...interface{}) {}

func (gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	slog.WarnContext(ctx, "gorm: "+fmt.Sprintf(msg, args...))
}

func (gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	slog.ErrorContext(ctx, "gorm: "+fmt.Sprintf(msg, args...))
}

// Trace se llama después de cada consulta SQL.
func (gormLogger) Trace(ctx context.Context, inicio time.Time, fc func() (string, int64), err error) {
	duracionMs := float64(time.Since(inicio).Microseconds()) / 1000

	switch {
	case err != nil && !errors.Is(err, logger.ErrRecordNotFound):
		sql, _ := fc()
		slog.ErrorContext(ctx, "error en una consulta SQL", "error", err, "sql", sql, "duration_ms", duracionMs)
	case duracionMs > 200:
		sql, _ := fc()
		slog.WarnContext(ctx, "consulta SQL lenta", "sql", sql, "duration_ms", duracionMs)
	}
}

// ParamsFilter hace que GORM nos pase la consulta sin los valores (con $1, $2...).
func (gormLogger) ParamsFilter(_ context.Context, sql string, _ ...interface{}) (string, []interface{}) {
	return sql, nil
}
