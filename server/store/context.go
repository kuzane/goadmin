package store

import (
    "context"

    "github.com/gin-gonic/gin"
)

const key = "store"

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Store {
    store, _ := c.Value(key).(Store)
    return store
}

func ToContext(c *gin.Context, store Store) {
    c.Set(key, store)
}
