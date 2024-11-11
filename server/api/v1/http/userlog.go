package http

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/kuzane/goadmin/server/model"
    "github.com/kuzane/goadmin/server/router/middleware/session"
    "github.com/kuzane/goadmin/server/store"
)

// UserlogDelete
//
//	@Tags       系统日志
//	@Summary    删除日志
//	@Router     /sys/logs/{id} [delete]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      id     path     string      true    "角色id"
func DeleteUserlog(ctx *gin.Context) {
    ids := session.Param(ctx)
    _store := store.FromContext(ctx)

    for _, id := range ids {
        log, err := _store.GetUserlogDetail(ctx, model.NewDescribeRequestByID(id))
        if err != nil {
            ctx.String(http.StatusBadRequest, "Error parsing Role id. %s", err.Error())
            return
        }

        if err := _store.DeleteUserlog(ctx, log); err != nil {
            ctx.String(http.StatusBadRequest, err.Error())
            return
        }
    }

    ctx.Status(http.StatusNoContent)
}

// UserlogList
//
//	@Tags       系统日志
//	@Summary    日志列表
//	@Router     /sys/logs [get]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    plain
//	@Success    200 {array}     model.Userlog  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      page            query   int     false   "for response pagination, page offset number"   default(1)
//	@Param      perPage         query   int     false   "for response pagination, max items per page"   default(50)
//	@Param      keyword         query   string      false   "根据关键字进行查询"
func GetUserlogList(ctx *gin.Context) {
    data, err := store.FromContext(ctx).GetUserlogList(ctx, &model.UserlogListOptions{
        ListOptions: session.Pagination(ctx),
        Keyword:     ctx.Query("keyword"),
    })

    if err != nil {
        ctx.String(http.StatusInternalServerError, "Error getting Role list. %s", err.Error())
        return
    }

    ctx.JSON(http.StatusOK, data)
}

// UserlogEmpty
//
//	@Tags       系统日志
//	@Summary    清空日志
//	@Router     /sys/logs [delete]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
func EmptyUserlog(ctx *gin.Context) {
    err := store.FromContext(ctx).EmptyUserlog(ctx)
    if err != nil {
        ctx.String(http.StatusInternalServerError, "清空日志失败!", err.Error())
        return
    }

    ctx.Status(http.StatusOK)
}
