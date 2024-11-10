package middleware

import (
    "github.com/gin-gonic/gin"

    "goamin/server/store"
)

func Store(v store.Store) gin.HandlerFunc {
    return func(c *gin.Context) {
        store.ToContext(c, v)
        c.Next()
    }
}
