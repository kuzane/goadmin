package session

import (
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/mssola/user_agent"
    "github.com/rs/zerolog/log"

    "github.com/kuzane/goadmin/server"
    "github.com/kuzane/goadmin/server/model"
    "github.com/kuzane/goadmin/server/store"
)

// SetUserlog 记录用户操作的中间件
func SetUserlog() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        start := time.Now()

        // 处理请求
        ctx.Next()

        // 请求结束后记录日志
        user := &model.User{}
        fullpath := ctx.FullPath()

        if fullpath == server.Config.Server.LoginPath {
            //如果是login接口，login的时候不能从上下文获取到user
            in := &model.Login{}
            ctx.ShouldBindJSON(in)
            user.Username = in.Username
        } else {
            user = User(ctx)
        }

        ua := user_agent.New(ctx.Request.UserAgent())
        browserName, browserVersion := ua.Browser()

        userlog := &model.Userlog{
            Username: user.Username,
            IPAddr:   ctx.ClientIP(),
            StartAt:  start.Unix(),
            Path:     strings.TrimPrefix(fullpath, server.Config.Server.RootPath),
            Method:   ctx.Request.Method,
            Status:   int64(ctx.Writer.Status()),
            Duration: time.Since(start).Milliseconds(),
            ClientOS: ua.OS(),
            Browser:  browserName + " " + browserVersion,
        }
        if err := store.FromContext(ctx).CreateUserlog(ctx, userlog); err != nil {
            // 创建失败打印日志记录
            log.Info().Msgf("创建userlog失败: %v", userlog)
        }
    }
}
