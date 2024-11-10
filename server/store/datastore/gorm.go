package datastore

import (
    "context"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gorm.io/gorm/logger"
)

type gormLogger struct {
    logger  zerolog.Logger
    level   logger.LogLevel
    showSQL bool
}

func newGORMLogger(level logger.LogLevel) *gormLogger {
    return &gormLogger{
        logger: log.With().Str("component", "gorm").Logger(),
        level:  level,
    }
}

func (g *gormLogger) LogMode(logger.LogLevel) logger.Interface {
    return g
}

func (g *gormLogger) Info(ctx context.Context, s string, v ...interface{}) {
    if g.level >= logger.Info {
        g.logger.Info().Msgf(s, v...)
    }
}

func (g *gormLogger) Warn(ctx context.Context, s string, v ...interface{}) {
    if g.level >= logger.Warn {
        g.logger.Warn().Msgf(s, v...)
    }
}

func (g *gormLogger) Error(ctx context.Context, s string, v ...interface{}) {
    if g.level >= logger.Error {
        g.logger.Warn().Msgf(s, v...)
    }
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    if g.level >= logger.Info {
        elapsed := time.Since(begin)
        sql, rows := fc()

        if err != nil {
            g.logger.Error().
                Err(err).
                Dur("elapsed", elapsed).
                Int64("rows", rows).
                Msgf("SQL: %s", sql)
        } else {
            g.logger.Info().
                Dur("elapsed", elapsed).
                Int64("rows", rows).
                Msgf("SQL: %s", sql)
        }
    }
}
