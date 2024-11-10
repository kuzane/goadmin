package app

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"

	"goamin/server"
	"goamin/server/store"
	"goamin/server/store/datastore"
)

func setupServerHost(c *cli.Command) error {
	if c.String("server-host") == "" {
		return fmt.Errorf("SERVER_HOST is not properly configured")
	}

	if !strings.Contains(c.String("server-host"), "://") {
		return fmt.Errorf("SERVER_HOST must be <scheme>://<hostname> format")
	}

	if _, err := url.Parse(c.String("server-host")); err != nil {
		return fmt.Errorf("could not parse SERVER_HOST: %w", err)
	}

	if strings.Contains(c.String("server-host"), "://localhost") {
		log.Warn().Msg(
			"SERVER_HOST should probably be publicly accessible (not localhost)",
		)
	}

	return nil
}

func setupStore(ctx context.Context, c *cli.Command) (store.Store, error) {
	datasource := c.String("datasource")
	driver := c.String("driver")
	maxOpenConn := c.Int("max-open-conn")
	maxIdleConn := c.Int("max-idle-conn")
	maxLifeTime := c.Int("max-life-time")
	maxIdleTime := c.Int("max-idle-time")
	isShowSql := c.Bool("log-orm-sql")

	if !datastore.SupportedDriver(driver) {
		return nil, fmt.Errorf("database driver '%s' not supported", driver)
	}

	opts := &store.Opts{
		Driver:      driver,
		Config:      datasource,
		MaxOpenConn: int(maxOpenConn),
		MaxIdleConn: int(maxIdleConn),
		MaxLifeTime: int(maxLifeTime),
		MaxIdleTime: int(maxIdleTime),
		ShowSQL:     isShowSql,
	}

	//配置gorm
	store, err := datastore.NewEngine(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := store.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("could not migrate datastore: %w", err)
	}

	return store, nil
}

func setupEvilGlobals(ctx context.Context, c *cli.Command, s store.Store) (err error) {
	err = setupSMTPConfig(c)

	serverHost := strings.TrimSuffix(c.String("server-host"), "/")
	server.Config.Server.JWTSecret, err = setupJWTSecret(ctx, s)
	server.Config.Server.RootUser = c.String("server-root-user")
	server.Config.Server.Host = serverHost
	server.Config.Server.RBAC, err = s.SetRBAC(ctx)
	server.Config.Server.RootPath = c.String("server-root-path")
	server.Config.Server.LoginPath = c.String("server-login-path")
	server.Config.Server.AccessTokenDuration = c.Int("server-access-atoken-duration")
	server.Config.Server.RefreshTokenDuration = c.Int("server-refresh-atoken-duration")
	server.Config.Server.EmailCaptchaTemp, err = setupSendEmailCaptcha(ctx, s)
	server.Config.Server.EmailPasswordTemp, err = setupSendEmailPassword(ctx, s)

	server.Config.Server.LoadWeb = c.Bool("load-web")
	server.Config.Server.ServerSwagger = c.Bool("server-swagger")

	return err
}

func setupSMTPConfig(c *cli.Command) (err error) {
	dsn := c.String("smtp-dsn")
	parts := strings.Split(dsn, ":") // 使用 ':' 分隔字符串
	if len(parts) != 5 {
		return fmt.Errorf("smtp-dsn配置错误")
	}
	server.Config.SmtpDSN.User = parts[0]
	server.Config.SmtpDSN.Pass = parts[1]
	server.Config.SmtpDSN.Host = parts[2]
	server.Config.SmtpDSN.Port, err = strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("smtp-dsn端口配置错误: %s", parts[3])
	}
	server.Config.SmtpDSN.Status, err = strconv.ParseBool(parts[4])
	if err != nil {
		return fmt.Errorf("smtp-dsn状态配置错误: %s", parts[4])
	}

	return nil
}

const (
	jwtSecretID     = "jwt-secret"
	emailCaptchaID  = "email-captcha"
	emailPasswordID = "email-password"
)

func setupJWTSecret(ctx context.Context, _store store.Store) (string, error) {
	jwtSecret, err := _store.ServerConfigGet(ctx, jwtSecretID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		jwtSecret := base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(32),
		)
		err = _store.ServerConfigSet(ctx, jwtSecretID, jwtSecret)
		if err != nil {
			return "", err
		}
		return jwtSecret, nil
	}
	if err != nil {
		return "", err
	}

	return jwtSecret, nil
}

func setupSendEmailCaptcha(ctx context.Context, _store store.Store) (string, error) {
	emailTemplate, err := _store.ServerConfigGet(ctx, emailCaptchaID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		emailTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>验证码邮件</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            padding: 20px;
            margin: 0;
        }
        .container {
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
        }
        h2 {
            color: #333333;
        }
        .code {
            font-size: 32px;
            font-weight: bold;
            color: #4CAF50;
            margin: 20px 0;
        }
        .note {
            font-size: 14px;
            color: #888888;
        }
        .footer {
            font-size: 12px;
            color: #888888;
            margin-top: 20px;
            text-align: center;
        }
        .footer a {
            color: #4CAF50;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>您好，{{.Name}}！</h2>
        <p>感谢您请求验证码。请使用以下验证码重置您的账户密码：</p>
        
        <div class="code">{{.Content}}</div>
        
        <p class="note">该验证码有效期为 5 分钟，请尽快使用。</p>

        <p class="note">如果您没有请求此验证码，请忽略此邮件。</p>
        
        <div class="footer">
            <p><a rel="noopener" href="{{.Domain}}" target="_blank">访问goAdmin</a></p>
        </div>
    </div>
</body>
</html>
`
		err = _store.ServerConfigSet(ctx, emailCaptchaID, emailTemplate)
		if err != nil {
			return "", err
		}
	}

	return emailTemplate, nil
}

func setupSendEmailPassword(ctx context.Context, _store store.Store) (string, error) {
	emailTemplate, err := _store.ServerConfigGet(ctx, emailPasswordID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		emailTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>您的账号密码</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            padding: 20px;
            margin: 0;
        }
        .container {
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
        }
        h2 {
            color: #333333;
        }
        .password {
            font-size: 24px;
            font-weight: bold;
            color: #ff5733;
            background-color: #f8f8f8;
            padding: 10px;
            border-radius: 5px;
            margin: 20px 0;
            word-wrap: break-word;
        }
        .note {
            font-size: 14px;
            color: #888888;
        }
        .footer {
            font-size: 12px;
            color: #888888;
            margin-top: 20px;
            text-align: center;
        }
        .footer a {
            color: #4CAF50;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>您好, {{.Name}}！</h2>
        <p>您的账号密码为：</p>
        
        <div class="password">{{.Content}}</div>
        
        <p class="note">为了确保您的账户安全，强烈建议您尽快修改密码。您可以在登录后访问“个人中心”进行修改。</p>

        <p class="note">如果您没有请求此邮件，请忽略此邮件。</p>
        
        <div class="footer">
            <p><a rel="noopener" href="{{.Domain}}" target="_blank">访问goAdmin</a></p>
        </div>
    </div>
</body>
</html>
`
		err = _store.ServerConfigSet(ctx, emailPasswordID, emailTemplate)
		if err != nil {
			return "", err
		}
	}

	return emailTemplate, nil
}
