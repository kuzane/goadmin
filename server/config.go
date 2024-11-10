package server

import (
	"github.com/mikespook/gorbac"
)

var Config = struct {
	// 服务相关的配置
	Server struct {
		Host                 string
		Port                 string
		JWTSecret            string
		RootPath             string
		RootUser             string
		LoginPath            string
		EmailCaptchaTemp     string
		EmailPasswordTemp    string
		AccessTokenDuration  int64
		RefreshTokenDuration int64
		LoadWeb              bool // 是否加载前端
		ServerSwagger        bool // 是否启用swagger
		RBAC                 *gorbac.RBAC
	}
	SmtpDSN struct {
		Host   string
		Port   int
		User   string
		Pass   string
		Status bool
	}
}{}
