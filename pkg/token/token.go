package token

import (
	"fmt"
	"goamin/server"
	"goamin/server/model"
	"goamin/server/store"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	claims jwt.MapClaims
}

func New() *Token {
	return &Token{claims: jwt.MapClaims{}}
}

func (t *Token) SignExpires(secret string, exp int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not a MapClaims")
	}

	for k, v := range t.claims {
		claims[k] = v
	}

	if exp > 0 {
		claims["exp"] = float64(exp)
	}

	return token.SignedString([]byte(secret))
}

func (t *Token) Get(key string) string {
	claim, ok := t.claims[key].(string)
	if !ok {
		return ""
	}

	return claim
}

func (t *Token) Set(key, value string) {
	t.claims[key] = value
}

// 解析token是否过期
func ParseToken(tokenString string, secret string) (*Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非法token: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		} else {
			return nil, fmt.Errorf("exp claim is missing")
		}

		return &Token{claims: claims}, nil
	}

	return nil, fmt.Errorf("token无效或解析失败")
}

func ParseRequest(ctx *gin.Context) (string, error) {
	// 从Header中获取Authorization字段，格式为 "Bearer {token}"
	ak := ctx.GetHeader("Authorization")
	if ak == "" || !strings.HasPrefix(ak, "Bearer ") {

		return "", fmt.Errorf("非法token！")
	}

	// 去掉前缀 "Bearer "
	ak = ak[7:]

	return ak, nil
}

func IsRefresh(t *Token) bool {
	expTime := int64(t.claims["exp"].(float64))
	currentTime := time.Now().Unix()
	if expTime-currentTime < 5*60 { // 如果距离过期小于5分钟
		return true
	}

	return false
}

// 重新签发ak
func SignAK(ctx *gin.Context, user *model.User, secret string) *model.User {
	rtk, err := ParseToken(user.RefreshToken, secret)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return user
	}

	// 重新签发AccessToken
	akexp := time.Now().Add(time.Duration(server.Config.Server.AccessTokenDuration) * time.Second).Unix()
	ak, err := rtk.SignExpires(secret, akexp)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return user
	}
	user.AccessToken = ak
	if err := store.FromContext(ctx).UpdateUser(ctx, user); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return user
	}

	ctx.Header("Access-Control-Expose-Headers", "token")
	ctx.Header("token", ak)

	return user
}
