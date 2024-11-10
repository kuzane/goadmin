package http

import (
	"errors"
	"goamin/server"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"goamin/pkg/utils"
	"goamin/server/model"
	"goamin/server/router/middleware/session"
	"goamin/server/store"
)

// RoleCreate
//
//	@Tags       系统角色
//	@Summary    创建角色
//	@Router     /sys/roles [post]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200 {object}    model.Role  "请求成功"
//	@Failure    400 {object}    string      "请求错误"
//	@Failure    500 {object}    string      "内部错误"
//	@Param      data            body    model.Role true "角色数据"
func PostRole(ctx *gin.Context) {
	in := new(model.CreateRole)

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}

	_store := store.FromContext(ctx)

	if err := _store.CreateRole(ctx, model.NewRole(in)); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, in)
}

// RoleDelete
//
//	@Tags       系统角色
//	@Summary    删除角色
//	@Router     /sys/roles/{id} [delete]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      id     path     string      true    "角色id"
func DeleteRole(ctx *gin.Context) {
	ids := session.Param(ctx)
	_store := store.FromContext(ctx)

	for _, id := range ids {
		role, err := _store.GetRoleDetail(ctx, model.NewDescribeRequestByID(id))
		if err != nil {
			ctx.String(http.StatusBadRequest, "查询角色失败", err.Error())
			return
		}

		if strings.ToLower(role.Rolename) == server.Config.Server.RootUser {
			ctx.String(http.StatusBadRequest, errors.New("不允许删除admin角色").Error())
			return
		}

		if err := _store.DeleteRole(ctx, role); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}

// RoleUpdate
//
//	@Tags       系统角色
//	@Summary    修改角色
//	@Router     /sys/roles/{id} [patch]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200 {object}    model.Role "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      id              path    string      true    "模板id"
//	@Param      data            body    model.Role  true    "更新的数据"
func PatchRole(ctx *gin.Context) {
	in := new(model.CreateRole)
	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}
	_store := store.FromContext(ctx)
	role, err := _store.GetRoleDetail(ctx, model.NewDescribeRequestByID(ctx.Param("id")))
	if err != nil {
		ctx.String(http.StatusBadRequest, utils.ErrRespString(err))
		return
	}

	if err := store.FromContext(ctx).UpdateRole(ctx, role.SetRole(in)); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// RolesList
//
//	@Tags       系统角色
//	@Summary    角色列表
//	@Router     /sys/roles [get]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    plain
//	@Success    200 {array}     model.Role  "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      page            query   int     false   "for response pagination, page offset number"   default(1)
//	@Param      perPage         query   int     false   "for response pagination, max items per page"   default(50)
//	@Param      keyword         query   string      false   "根据关键字进行查询"
//	@Param      rolename        query   string      false   "根据Rolename进行查询"
//	@Param      nickname        query   string      false   "根据nickname进行查询"
func GetRoleList(ctx *gin.Context) {
	data, err := store.FromContext(ctx).GetRoleList(ctx, &model.RoleListOptions{
		ListOptions: session.Pagination(ctx),
		Keyword:     ctx.Query("keyword"),
		Rolename:    ctx.Query("rolename"),
		Nickname:    ctx.Query("nickname"),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error getting Role list. %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// RoleDetails
//
//	@Tags       系统角色
//	@Summary    角色详情
//	@Router     /sys/roles/{id} [get]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200 {object}    model.Role "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
//	@Param      id          path    string      true    "角色id"
func GetRoleDetails(ctx *gin.Context) {
	Role, err := store.FromContext(ctx).GetRoleDetail(ctx, model.NewDescribeRequestByID(ctx.Param("id")))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Error parsing Role id. %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Role)
}

// EndpointTree
//
//	@Tags       系统角色
//	@Summary    获取接口树
//	@Router     /sys/roles/apis [get]
//	@Param      Authorization   header  string  true    "Insert your personal access token" default(Bearer <personal access token>)
//	@Produce    json
//	@Success    200 {object}    model.Endpoint "成功"
//	@Failure    400 {object}    string "请求错误"
//	@Failure    500 {object}    string "内部错误"
func GetEndpointTree(ctx *gin.Context) {
	data, err := store.FromContext(ctx).GetEndpointList(ctx, model.NewEndpointListAll())
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error getting Endpoint list. %s", err.Error())
		return
	}

	tree := model.GenerateTree(data.Items)
	ctx.JSON(http.StatusOK, tree)
}
