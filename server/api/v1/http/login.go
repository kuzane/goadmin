package http

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"goamin/pkg/token"
	"goamin/pkg/utils"
	"goamin/server"
	"goamin/server/model"
	"goamin/server/store"
)

var (
	captcha = cache.New(5*time.Minute, 10*time.Minute)
)

// Login
//
//	@Tags		必开接口
//	@Summary	用户登录
//	@Router		/login [post]
//	@Produce	json
//	@Success	200	{object}    string "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      data            body    model.Login true "登录"
func Login(ctx *gin.Context) {
	in := &model.Login{}
	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}
	_store := store.FromContext(ctx)

	user, err := _store.GetUserDetail(ctx, model.NewNewDetailUserRequestByName(in.Username))
	if err != nil {
		ctx.String(http.StatusBadRequest, "账户错误！")
		return
	}
	if !user.Status {
		ctx.String(http.StatusBadRequest, "账户已被禁用！")
		return
	}

	if err := user.PasswordVerifiers(in.Password); err != nil {
		ctx.String(http.StatusBadRequest, "密码错误！")
		return
	}

	//颁发token
	tk := token.New()
	secret := server.Config.Server.JWTSecret
	akexp := time.Now().Add(time.Duration(server.Config.Server.AccessTokenDuration) * time.Second).Unix()
	ak, err := tk.SignExpires(secret, akexp)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "服务内部错误", err.Error())
		return
	}

	rkexp := time.Now().Add(time.Duration(server.Config.Server.RefreshTokenDuration) * time.Second).Unix()
	rk, err := tk.SignExpires(secret, rkexp)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "服务内部错误", err.Error())
		return
	}
	user.AccessToken = ak
	user.RefreshToken = rk
	_store.UpdateUser(ctx, user)

	ctx.JSON(http.StatusOK, ak)
}

// Logout
//
//	@Tags		必开接口
//	@Summary	退出登陆
//
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /logout [post]
// @Security
func Logout(ctx *gin.Context) {
	// 记录日志
	ctx.Status(http.StatusOK)
}

// GetCaptchaByEmail
//
//	@Tags		必开接口
//	@Summary	获取邮箱验证码
//	@Router		/email/captcha [post]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}    string "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
func GetCaptchaByEmail(ctx *gin.Context) {
	param := new(struct {
		Email string `json:"email" binding:"required"`
	})
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}

	// 1.邮箱验证(通过邮箱能否查到用户)
	user, err := store.FromContext(ctx).GetUserDetail(ctx, model.NewNewDetailUserRequestByEmail(param.Email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.String(http.StatusBadRequest, "根据邮箱未能查到您的账户，请您先注册！")
			return
		}
		ctx.String(http.StatusInternalServerError, "系统错误，请联系管理员！", err.Error())
		return
	}

	// 2.生成随机验证码
	code := utils.GenerateCaptcha()
	captcha.Set(param.Email, code, cache.DefaultExpiration)

	// 3.邮箱发送，启用邮箱就发送
	if server.Config.SmtpDSN.Status {
		if err = utils.SendEmail(
			server.Config.SmtpDSN.Host,
			server.Config.SmtpDSN.Port,
			server.Config.SmtpDSN.User,
			server.Config.SmtpDSN.Pass,
			[]string{param.Email},
			[]string{},
			"「重置密码」验证码",
			server.Config.Server.EmailCaptchaTemp,
			utils.EmailBodyData{
				Name:    user.Nickname,
				Content: code,
				Domain:  server.Config.Server.Host,
			}); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	ctx.Status(http.StatusOK)
}

// ResetPwdByEmail
//
//	@Tags		必开接口
//	@Summary	通过邮箱重置密码
//	@Router		/pwd/forgot [post]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}    string "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
func ResetPwdByEmail(ctx *gin.Context) {
	param := new(struct {
		Email   string `json:"email" validate:"required"`
		Captcha string `json:"captcha" validate:"required"`
	})
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}

	// 1.验证码校验
	get, found := captcha.Get(param.Email)
	if !found {
		ctx.String(http.StatusNotFound, "验证码已过期")
		return
	}
	if get.(string) != param.Captcha {
		ctx.String(http.StatusNotFound, "输入验证码错误")
		return
	}

	// 2.邮箱验证(通过邮箱能否查到用户)
	_store := store.FromContext(ctx)
	user, err := _store.GetUserDetail(ctx, model.NewNewDetailUserRequestByEmail(param.Email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.String(http.StatusBadRequest, "根据邮箱未能查到您的账户，请您先注册！")
			return
		}
		ctx.String(http.StatusInternalServerError, "系统错误，请联系管理员！", err.Error())
		return
	}

	// 3.生成随机密码
	pass := utils.GeneratePassword()

	// 4.修改用户密码
	user.Password = pass
	if err := _store.UpdateUser(ctx, user); err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误，请联系管理员！", err.Error())
		return
	}

	// 5.邮件通知
	if server.Config.SmtpDSN.Status {
		if err = utils.SendEmail(
			server.Config.SmtpDSN.Host,
			server.Config.SmtpDSN.Port,
			server.Config.SmtpDSN.User,
			server.Config.SmtpDSN.Pass,
			[]string{param.Email},
			[]string{},
			"「重置密码」密码",
			server.Config.Server.EmailPasswordTemp,
			utils.EmailBodyData{
				Name:    user.Nickname,
				Content: pass,
				Domain:  server.Config.Server.Host,
			}); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	ctx.Status(http.StatusOK)
}
