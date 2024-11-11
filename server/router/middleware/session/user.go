package session

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/kuzane/goadmin/pkg/token"
    "github.com/kuzane/goadmin/server"
    "github.com/kuzane/goadmin/server/model"
    "github.com/kuzane/goadmin/server/store"
)

func User(c *gin.Context) *model.User {
    v, ok := c.Get("user")
    if !ok {
        return nil
    }
    u, ok := v.(*model.User)
    if !ok {
        return nil
    }
    return u
}

func SetUser() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var user *model.User
        secret := server.Config.Server.JWTSecret
        // 1.获取ak
        ak, err := token.ParseRequest(ctx)
        if err != nil || ak == "" {
            ctx.String(http.StatusUnauthorized, err.Error())
            ctx.Abort()
            return
        }

        // 2.通过ak获取用户
        if user, err = store.FromContext(ctx).GetUserDetail(ctx, model.NewNewDetailUserRequestByAK(ak)); err != nil {
            if err == gorm.ErrRecordNotFound {
                ctx.String(http.StatusUnauthorized, "该账户已在其他地方登录?")
                ctx.Abort()
                return
            }
            ctx.String(http.StatusInternalServerError, err.Error())
            ctx.Abort()
            return
        }
        // 2.验证是否合法
        _, err = token.ParseToken(ak, secret)
        if err != nil {
            if err.Error() == "token has invalid claims: token is expired" {
                //验证ak是否过期，若ak过期，数据库中rk有效则颁发新token
                user = token.SignAK(ctx, user, secret)
                ctx.Set("user", user)
                ctx.Next()
            }
            ctx.String(http.StatusUnauthorized, err.Error())
            ctx.Abort()
            return
        }

        ctx.Set("user", user)
        ctx.Next()
    }
}
