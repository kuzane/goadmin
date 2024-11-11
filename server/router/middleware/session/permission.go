package session

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/mikespook/gorbac"

    "github.com/kuzane/goadmin/server"
)

func RoleCheckAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        user := User(ctx)
        // 对admin用户放行所有接口
        if strings.ToLower(user.Username) == server.Config.Server.RootUser {
            ctx.Next()
            return
        }

        method := ctx.Request.Method
        path := strings.TrimPrefix(ctx.FullPath(), server.Config.Server.RootPath)
        permission := method + ":" + path
        // 权限应该合并校验
        isPermission := false

        for _, role := range user.Roles {
            // 验证角色是否有访问接口的权限
            if server.Config.Server.RBAC.IsGranted(role.Rolename, gorbac.NewStdPermission(permission), nil) {
                isPermission = true
                break
            }
        }

        if isPermission {
            ctx.Next()
        } else {
            ctx.String(http.StatusForbidden, "您没有权限执行此操作，请您联系管理员进行授权！")
            ctx.Abort()
        }
    }
}
