package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

func Logger(timeFormat string, utc bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        // some evil middlewares modify this values
        path := c.Request.URL.Path
        c.Next()

        end := time.Now()
        latency := end.Sub(start)
        if utc {
            end = end.UTC()
        }

        entry := map[string]any{
            "status":     c.Writer.Status(),
            "method":     c.Request.Method,
            "path":       path,
            "ip":         c.ClientIP(),
            "latency":    latency,
            "user-agent": c.Request.UserAgent(),
            "time":       end.Format(timeFormat),
        }

        if len(c.Errors) > 0 {
            // Append error field if this is an erroneous request.
            log.Error().Str("error", c.Errors.String()).Fields(entry).Msg("")
        } else {
            log.Debug().Fields(entry).Msg("")
        }

    }
}
