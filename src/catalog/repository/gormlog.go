package repository

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

// gormSlog es un "adaptador": GORM sabe escribir en cualquier cosa que tenga
// un método Printf, así que le damos uno que manda el mensaje a slog (JSON).
type gormSlog struct{}

func (gormSlog) Printf(format string, args ...interface{}) {
	slog.Warn("gorm", "detalle", fmt.Sprintf(format, args...))
}

// newGormLogger configura qué cosas registra GORM:
// solo errores reales y consultas lentas (más de 200 ms),
// sin el "record not found" que aparece en cada 404.
func newGormLogger() logger.Interface {
	return logger.New(gormSlog{}, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
}
