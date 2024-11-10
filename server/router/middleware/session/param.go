package session

import (
    "strings"

    "github.com/gin-gonic/gin"
)

// gin获取路径参数,主要用于删除多个对象时
func Param(ctx *gin.Context) []string {
    ids := ctx.Param("id")
    idSlice := strings.Split(ids, ",")

    return idSlice
}
