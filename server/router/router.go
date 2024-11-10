package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"net/url"

	"goamin/docs"
	"goamin/server"
	"goamin/server/router/middleware/header"
	"goamin/server/store"
)

func Load(_store store.Store, middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	e.UseRawPath = true
	e.Use(gin.Recovery())
	e.Use(func(c *gin.Context) {
		log.Trace().Msgf("[%s] %s", c.Request.Method, c.Request.URL.String())
		c.Next()
	})

	if server.Config.Server.ServerSwagger {
		setupSwaggerConfigAndRoutes(e)
	}

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)
	base := e.Group("")
	{
		// 放置promethean指标接口

		if server.Config.Server.LoadWeb {
			// 前后端不分离时
			e.Static("/static", "./web/dist")
			e.GET("/", func(ctx *gin.Context) {
				ctx.File("./web/dist/index.html")
			})
			// 如果 Vue 使用前端路由，确保其他未匹配的路径都返回 `index.html`
			e.NoRoute(func(c *gin.Context) {
				c.File("./web/dist/index.html")
			})
		}
	}
	apiRoutes(base, _store)
	return e
}
func setupSwaggerConfigAndRoutes(e *gin.Engine) {
	docs.SwaggerInfo.Host = getHost(server.Config.Server.Host)
	docs.SwaggerInfo.BasePath = server.Config.Server.RootPath
	e.GET(""+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
func getHost(s string) string {
	parse, _ := url.Parse(s)
	return parse.Host
}
