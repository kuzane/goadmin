package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload" // 读取项目.env文件
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"goamin/pkg/logger"
	"goamin/pkg/utils"
	"goamin/server/app"
)

var flags = append([]cli.Flag{
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_HOST"),
		Name:    "server-host",
		Usage:   "服务访问的完全的URL！例如: <scheme>://<host>[/<prefix path>]",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_PORT"),
		Name:    "server-port",
		Usage:   "服务启动端口",
		Value:   "8000",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_ROOT_USER"),
		Name:    "server-root-user",
		Usage:   "系统管理用户名",
		Value:   "admin",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_ROOT_PATH"),
		Name:    "server-root-path",
		Usage:   "url的前缀路径",
		Value:   "/api/v1",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_LOGIN_PATH"),
		Name:    "server-login-path",
		Usage:   "登录接口的路径",
		Value:   "/login",
	},
	&cli.IntFlag{ //access_token
		Sources: cli.EnvVars("SERVER_ACCESS_TOKEN_DURATION"),
		Name:    "server-access-atoken-duration",
		Usage:   "access-atoken有效期时长(默认2h)",
		Value:   7200, // 7200
	},
	&cli.IntFlag{ //refresh_token
		Sources: cli.EnvVars("SERVER_REFRESH_TOKEN_DURATION"),
		Name:    "server-refresh-atoken-duration",
		Usage:   "refresh-atoken有效期时长(默认7天)",
		Value:   86400 * 7,
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SERVER_TRANS"),
		Name:    "server-trans",
		Usage:   "翻译器",
		Value:   "zh",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("SMTP_DSN"),
		Name:    "smtp-dsn",
		Usage:   "邮箱DSN(默认不启用)",
		Value:   "::::false",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("DATABASE_DRIVER"),
		Name:    "driver",
		Usage:   "数据库驱动",
		Value:   "mysql",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("DATABASE_DATASOURCE"),
		Name:    "datasource",
		Usage:   "数据库驱动配置(DSN)",
	},
	&cli.IntFlag{
		Sources: cli.EnvVars("DATABASE_MAX_OPEN_CONN"),
		Name:    "max-open-conn",
		Usage:   "数据库最大打开的连接数",
	},
	&cli.IntFlag{
		Sources: cli.EnvVars("DATABASE_MAX_IDLE_CONN"),
		Name:    "max-idle-conn",
		Usage:   "数据库连接池的最大空闲连接",
	},
	&cli.IntFlag{
		Sources: cli.EnvVars("DATABASE_MAX_LIFE_TIME"),
		Name:    "max-life-time",
		Usage:   "数据库连接池中连接的最大生命周期",
	},
	&cli.IntFlag{
		Sources: cli.EnvVars("DATABASE_MAX_IDLE_TIME"),
		Name:    "max-idle-time",
		Usage:   "数据库连接池中连接的最大空闲时间",
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("LOG_ORM_SQL"),
		Name:    "log-orm-sql",
		Usage:   "日志是否打印SQL",
		Value:   false,
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("LOAD_WEB"),
		Name:    "load-web",
		Usage:   "前后端分离",
		Value:   false,
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("SERVER_SWAGGER"),
		Name:    "server-swagger",
		Usage:   "是否生成接口文档",
		Value:   false,
	},
}, logger.GlobalLoggerFlags...)

func main() {
	ctx := utils.WithContextSigtermCallback(context.Background(), func() {
		log.Info().Msg("收到终止信号, 关闭服务")
	})

	cmd := cli.Command{
		Name:     "goadmin-serve",
		Usage:    "goadmin server",
		Version:  "v1.0.0",
		Flags:    flags,
		Action:   app.Run,
		Commands: []*cli.Command{},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Error().Err(err).Msgf("运行服务时报错")
	}
}
