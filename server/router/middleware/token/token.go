package token

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"goamin/pkg/token"
	"goamin/server"
	"goamin/server/router/middleware/session"
)

func Refresh(ctx *gin.Context) {
	user := session.User(ctx)
	secret := server.Config.Server.JWTSecret
	atk, err := token.ParseToken(user.AccessToken, secret)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if token.IsRefresh(atk) { // 1. 判断ak有效期是否小于5分钟
		token.SignAK(ctx, user, secret)
	}

	ctx.Next()
}
