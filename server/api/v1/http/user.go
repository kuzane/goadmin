package http

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/kuzane/goadmin/server"

    "github.com/gin-gonic/gin"

    "github.com/kuzane/goadmin/pkg/utils"
    "github.com/kuzane/goadmin/server/model"
    "github.com/kuzane/goadmin/server/router/middleware/session"
    "github.com/kuzane/goadmin/server/store"
)

// UserCreate
//
//	@Tags		系统用户
//	@Summary	创建用户
//	@Router		/sys/users [post]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success    200 {object}    model.User  "请求成功"
//	@Failure    400 {object}    string      "请求错误"
//	@Failure    500 {object}    string      "内部错误"
//	@Param      data            body    model.User true "用户数据"
func PostUser(ctx *gin.Context) {
    in := &model.CreateUser{}
    if err := ctx.ShouldBindJSON(in); err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }
    // 创建用户没有传密码就生成一个6位随机密码
    if in.Password == "" {
        in.Password = utils.GeneratePassword()
    }
    user, err := model.NewUser(in)
    if err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }

    if err := store.FromContext(ctx).CreateUser(ctx, user); err != nil {
        ctx.String(http.StatusInternalServerError, utils.ErrRespString(err))
        return
    }
    // 将随机密码通过邮件发送给用户
    if server.Config.SmtpDSN.Status {
        if err = utils.SendEmail(
            server.Config.SmtpDSN.Host,
            server.Config.SmtpDSN.Port,
            server.Config.SmtpDSN.User,
            server.Config.SmtpDSN.Pass,
            []string{user.Email},
            []string{},
            "「感谢注册」密码",
            server.Config.Server.EmailPasswordTemp,
            utils.EmailBodyData{
                Name:    user.Nickname,
                Content: in.Password,
                Domain:  server.Config.Server.Host,
            }); err != nil {
            ctx.String(http.StatusBadRequest, err.Error())
            return
        }
    }

    ctx.JSON(http.StatusOK, in)
}

// UserDelete
//
//	@Tags		系统用户
//	@Summary	删除用户
//	@Router		/sys/users/{id} [delete]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	 "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param		id	   path	    string	    true	"用户id"
func DeleteUser(ctx *gin.Context) {
    ids := session.Param(ctx)
    _store := store.FromContext(ctx)

    for _, id := range ids {
        user, err := _store.GetUserDetail(ctx, model.NewDetailUserRequestByID(id))
        if err != nil {
            ctx.String(http.StatusBadRequest, fmt.Sprintf("当id=%s时查询报错: %s", id, err.Error()))
            return
        }

        if strings.ToLower(user.Username) == server.Config.Server.RootUser {
            ctx.String(http.StatusBadRequest, fmt.Sprintf("不允许删除admin账户"))
            return
        }

        if err := _store.DeleteUser(ctx, user); err != nil {
            ctx.String(http.StatusBadRequest, err.Error())
            return
        }
    }

    ctx.Status(http.StatusNoContent)
}

// UserUpdate
//
//	@Tags		系统用户
//	@Summary	修改用户
//	@Router		/sys/users/{id} [patch]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}	model.User "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param		id		        path	string	    true    "模板id"
//	@Param      data            body    model.User  true    "更新的数据"
func PatchUser(ctx *gin.Context) {
    in := new(model.CreateUser)
    if err := ctx.ShouldBindJSON(in); err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }
    _store := store.FromContext(ctx)
    user, err := _store.GetUserDetail(ctx, model.NewDetailUserRequestByID(ctx.Param("id")))
    if err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }

    setUser, err := user.SetUser(in)
    if err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }

    if err := store.FromContext(ctx).UpdateUser(ctx, setUser); err != nil {
        ctx.String(http.StatusInternalServerError, "Error updating secret %q. %s", in.Username, err.Error())
        return
    }

    ctx.Status(http.StatusOK)
}

// UsersList
//
//	@Tags		系统用户
//	@Summary	用户列表
//	@Router		/sys/users [get]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	plain
//	@Success	200 {array}     model.User  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param	    page			query	int		false	"for response pagination, page offset number"	default(1)
//	@Param		perPage			query	int		false	"for response pagination, max items per page"	default(50)
//	@Param		keyword		    query	string	    false	"根据关键字进行查询"
//	@Param		username		query	string	    false	"根据username进行查询"
//	@Param		nickname		query	string	    false	"根据nickname进行查询"
//	@Param		email		    query	string	    false	"根据email进行查询"
//	@Param		phone		    query	string	    false	"根据phone进行查询"
func GetUserList(ctx *gin.Context) {
    data, err := store.FromContext(ctx).GetUserList(ctx, &model.UserListOptions{
        ListOptions: session.Pagination(ctx),
        Keyword:     ctx.Query("keyword"),
        Username:    ctx.Query("username"),
        Nickname:    ctx.Query("nickname"),
        Email:       ctx.Query("email"),
        Phone:       ctx.Query("phone"),
    })

    if err != nil {
        ctx.String(http.StatusInternalServerError, "Error getting user list. %s", err.Error())
        return
    }

    ctx.JSON(http.StatusOK, data)
}

// UserDetails
//
//	@Tags		系统用户
//	@Summary	用户详情
//	@Router		/sys/users/{id} [get]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}	model.User "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param		id			path	string	    true	"用户id"
func GetUserDetails(ctx *gin.Context) {
    user, err := store.FromContext(ctx).GetUserDetail(ctx, model.NewDetailUserRequestByID(ctx.Param("id")))
    if err != nil {
        ctx.String(http.StatusBadRequest, "Error parsing user id. %s", err.Error())
        return
    }

    ctx.JSON(http.StatusOK, user)
}

// UserChangePassword
//
//	@Tags		系统用户
//	@Summary	修改密码
//	@Router		/password [put]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}	model.ChangePasswordRequest "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param		id		        path	string	    true    "id"
//	@Param      data            body    model.ChangePasswordRequest  true    "更新的数据"
func ChangePasswordUser(ctx *gin.Context) {
    in := new(model.ChangePasswordRequest)
    if err := ctx.ShouldBindJSON(in); err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }
    _store := store.FromContext(ctx)
    user := session.User(ctx)

    u, err := user.SetUserPasswordByChangePasswordRequest(in)
    if err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }
    if err := _store.UpdateUser(ctx, u); err != nil {
        ctx.String(http.StatusInternalServerError, "Error updating secret %q. %s", u.Username, err.Error())
        return
    }

    ctx.Status(http.StatusOK)
}

// UserChangeProfile
//
//	@Tags		系统用户
//	@Summary	修改个人信息
//	@Router		/password [put]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}	model.ChangePasswordRequest "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param		id		        path	string	    true    "id"
//	@Param      data            body    model.ChangePasswordRequest  true    "更新的数据"
func ChangeProfileUser(ctx *gin.Context) {
    in := new(model.ChangeProfileRequest)
    if err := ctx.ShouldBindJSON(in); err != nil {
        ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
        return
    }
    user := session.User(ctx)

    if err := store.FromContext(ctx).UpdateUser(ctx, user.SetUserPasswordByChangeProfileRequest(in)); err != nil {
        ctx.String(http.StatusInternalServerError, "Error updating secret %q. %s", user.Username, err.Error())
        return
    }

    ctx.Status(http.StatusOK)
}

// UserInfo
//
//	@Tags		系统用户
//	@Summary	用户详情
//	@Router		/sys/users/info [get]
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Produce	json
//	@Success	200	{object}	model.User "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
func GetUserInfo(ctx *gin.Context) {
    user := session.User(ctx)
    data := new(struct {
        Name        string   `json:"name"` // username
        Avatar      string   `json:"avatar"`
        Roles       []string `json:"roles"`
        Permissions []string `json:"permissions"` //用户的接口权限
        Phone       string   `json:"phone"`
        Email       string   `json:"email"`
        Nickname    string   `json:"nickname"`
        Description string   `json:"description"`
    })
    for _, role := range user.Roles {
        data.Roles = append(data.Roles, role.Rolename)
        data.Roles = append(data.Roles, role.Parents...)
    }
    data.Roles = removeDuplicates(data.Roles)
    data.Name = user.Username
    data.Avatar = user.Avatar
    data.Email = user.Email
    data.Phone = user.Phone
    data.Description = user.Description
    data.Nickname = user.Nickname

    if user.Username == server.Config.Server.RootUser {
        data.Permissions = []string{"*"}
    } else {
        for _, name := range data.Roles {
            detail, err := store.FromContext(ctx).GetRoleDetail(ctx, model.NewDescribeRequestByName(name))
            if err != nil {
                ctx.String(http.StatusInternalServerError, err.Error())
                return
            }
            if detail.Status {
                for _, endpoint := range detail.Endpoints {
                    data.Permissions = append(data.Permissions, endpoint.Identity)
                }
            }
        }
        data.Permissions = removeDuplicates(data.Permissions)
    }

    ctx.JSON(http.StatusOK, data)
}

func removeDuplicates(strings []string) []string {
    // 使用 map 来记录出现过的元素
    seen := make(map[string]struct{})
    var result []string

    for _, str := range strings {
        if _, ok := seen[str]; !ok {
            // 如果没有见过这个元素，添加到结果中
            seen[str] = struct{}{}
            result = append(result, str)
        }
    }

    return result
}
