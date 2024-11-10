package store

import (
	"context"

	"github.com/mikespook/gorbac"

	"goamin/server/model"
)

type Store interface {
	// ServerConfig
	ServerConfigGet(c context.Context, k string) (v string, err error)
	ServerConfigSet(c context.Context, k string, v string) error
	ServerConfigDelete(c context.Context, k string) error

	// users
	CreateUser(c context.Context, v *model.User) error
	DeleteUser(c context.Context, v *model.User) error
	UpdateUser(c context.Context, v *model.User) error
	GetUserDetail(c context.Context, req *model.DetailUserRequest) (data *model.User, err error)
	GetUserList(c context.Context, req *model.UserListOptions) (data *model.UserSet, err error)
	GetUserCount(c context.Context) (count int64, err error)

	// roles
	CreateRole(c context.Context, v *model.Role) error
	DeleteRole(c context.Context, v *model.Role) error
	UpdateRole(c context.Context, v *model.Role) error
	GetRoleDetail(c context.Context, req *model.DetailOptions) (data *model.Role, err error)
	GetRoleList(c context.Context, req *model.RoleListOptions) (data *model.RoleSet, err error)
	GetRoleCount(c context.Context) (count int64, err error)

	// endpoint
	CreateEndpoint(c context.Context, v *model.Endpoint) error
	DeleteEndpoint(c context.Context, v *model.Endpoint) error
	UpdateEndpoint(c context.Context, v *model.Endpoint) error
	GetEndpointDetail(c context.Context, req *model.DetailOptions) (data *model.Endpoint, err error)
	GetEndpointList(c context.Context, req *model.EndpointListOptions) (data *model.EndpointSet, err error)
	GetEndpointCount(c context.Context) (count int64, err error)
	SetEndpoint(c context.Context, v *model.Endpoint) error

	// userlog
	CreateUserlog(c context.Context, v *model.Userlog) error
	DeleteUserlog(c context.Context, v *model.Userlog) error
	GetUserlogDetail(c context.Context, req *model.DetailOptions) (data *model.Userlog, err error)
	GetUserlogList(c context.Context, req *model.UserlogListOptions) (data *model.UserlogSet, err error)
	EmptyUserlog(c context.Context) error

	// RBAC
	SetRBAC(c context.Context) (*gorbac.RBAC, error)

	Migrate(context.Context) error
	Close() error
}
