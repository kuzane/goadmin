package router

import (
    "context"
    "fmt"
    "os"
    "strings"

    "goamin/server/router/middleware/token"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"

    "goamin/server"
    "goamin/server/api/v1/http"
    "goamin/server/model"
    "goamin/server/router/middleware/session"
    "goamin/server/store"
)

func apiRoutes(e *gin.RouterGroup, _store store.Store) {
    apiBase := e.Group(server.Config.Server.RootPath)

    {
        { // 无需认证和监权
            apiBase.POST(server.Config.Server.LoginPath, http.Login)
            apiBase.POST("/email/captcha", http.GetCaptchaByEmail)
            apiBase.POST("/pwd/forgot", http.ResetPwdByEmail)
        }

        // 认证中间件
        apiBase.Use(session.SetUser())
        // access_token 刷新
        apiBase.Use(token.Refresh)
        // 用户操作中间件
        apiBase.Use(session.SetUserlog())
        {
            // 仅需要认证的接口，必开接口
            apiBase.POST("/logout", http.Logout)
            apiBase.GET("/userinfo", http.GetUserInfo)
            apiBase.PUT("/password", http.ChangePasswordUser)
            apiBase.PUT("/profile", http.ChangeProfileUser)
        }

        // 监权中间件
        apiBase.Use(session.RoleCheckAuth())
        // 系统相关接口
        sys := apiBase.Group("/sys")
        {
            users := sys.Group("/users")
            {
                registerRoute(users, &model.CreateEndpoint{Path: "", Method: "POST", Module: "系统管理", Kind: "用户管理", Identity: "add_user", Remark: "增加用户"}, http.PostUser, _store)
                registerRoute(users, &model.CreateEndpoint{Path: "/:id", Method: "DELETE", Module: "系统管理", Kind: "用户管理", Identity: "del_user", Remark: "删除用户"}, http.DeleteUser, _store)
                registerRoute(users, &model.CreateEndpoint{Path: "/:id", Method: "PATCH", Module: "系统管理", Kind: "用户管理", Identity: "upd_user", Remark: "修改用户"}, http.PatchUser, _store)
                registerRoute(users, &model.CreateEndpoint{Path: "/:id", Method: "GET", Module: "系统管理", Kind: "用户管理", Identity: "get_user", Remark: "用户详情"}, http.GetUserDetails, _store)
                registerRoute(users, &model.CreateEndpoint{Path: "", Method: "GET", Module: "系统管理", Kind: "用户管理", Identity: "list_user", Remark: "用户列表"}, http.GetUserList, _store)
            }

            roles := sys.Group("/roles")
            {
                registerRoute(roles, &model.CreateEndpoint{Path: "", Method: "POST", Module: "系统管理", Kind: "角色管理", Identity: "add_role", Remark: "增加角色"}, http.PostRole, _store)
                registerRoute(roles, &model.CreateEndpoint{Path: "/:id", Method: "DELETE", Module: "系统管理", Kind: "角色管理", Identity: "del_role", Remark: "删除角色"}, http.DeleteRole, _store)
                registerRoute(roles, &model.CreateEndpoint{Path: "/:id", Method: "PATCH", Module: "系统管理", Kind: "角色管理", Identity: "upd_role", Remark: "修改角色"}, http.PatchRole, _store)
                registerRoute(roles, &model.CreateEndpoint{Path: "/:id", Method: "GET", Module: "系统管理", Kind: "角色管理", Identity: "get_role", Remark: "角色详情"}, http.GetRoleDetails, _store)
                registerRoute(roles, &model.CreateEndpoint{Path: "", Method: "GET", Module: "系统管理", Kind: "角色管理", Identity: "list_role", Remark: "角色列表"}, http.GetRoleList, _store)
                registerRoute(roles, &model.CreateEndpoint{Path: "/apis", Method: "GET", Module: "系统管理", Kind: "角色管理", Identity: "tree_api", Remark: "接口树"}, http.GetEndpointTree, _store)
            }

            userlog := sys.Group("/logs")
            {
                registerRoute(userlog, &model.CreateEndpoint{Path: "/:id", Method: "DELETE", Module: "系统管理", Kind: "日志管理", Identity: "del_log", Remark: "删除日志"}, http.DeleteUserlog, _store)
                registerRoute(userlog, &model.CreateEndpoint{Path: "", Method: "GET", Module: "系统管理", Kind: "日志管理", Identity: "list_log", Remark: "日志列表"}, http.GetUserlogList, _store)
                registerRoute(userlog, &model.CreateEndpoint{Path: "", Method: "DELETE", Module: "系统管理", Kind: "日志管理", Identity: "reset_log", Remark: "清空日志"}, http.EmptyUserlog, _store)
            }
        }
    }
}

// registerRoute注册路由并将当前路由自动入库,后续实现相关RBAC
func registerRoute(group *gin.RouterGroup, v *model.CreateEndpoint, handler gin.HandlerFunc, _store store.Store) {
    router_path := v.Path
    v.Path = strings.TrimPrefix(fmt.Sprintf("%s%s", group.BasePath(), v.Path), server.Config.Server.RootPath)
    // 自动处理接口的增删改
    if err := _store.SetEndpoint(context.Background(), model.NewEndpoint(v)); err != nil {
        log.Debug().Msgf("Failed to log route: %v", err)
        os.Exit(100)
    }

    // 注册实际的路由
    switch v.Method {
    case "GET":
        group.GET(router_path, handler)
    case "POST":
        group.POST(router_path, handler)
    case "PUT":
        group.PUT(router_path, handler)
    case "PATCH":
        group.PATCH(router_path, handler)
    case "DELETE":
        group.DELETE(router_path, handler)
    default:
        log.Printf("Unsupported method: %s", router_path)
    }
}
