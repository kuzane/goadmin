package session

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "goamin/server/model"
)

const maxPageSize = 50

func Pagination(ctx *gin.Context) *model.ListOptions {
    page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
    if err != nil || page < 1 {
        page = 1
    }
    perPage, err := strconv.ParseInt(ctx.Query("perPage"), 10, 64)
    if err != nil || perPage < 1 || perPage > maxPageSize {
        perPage = maxPageSize
    }
    return &model.ListOptions{
        Page:    int(page),
        PerPage: int(perPage),
    }
}
