package datastore_test

import (
    "testing"

    "github.com/kuzane/goadmin/server/model"
)

func TestCreateRole(t *testing.T) {
    err := engine.CreateRole(ctx, &model.Role{
        Rolename:    "admin",
        Nickname:    "管理员",
        Status:      true,
        Description: "管理员权限",
        Parents:     nil,
        Users: []*model.User{
            {Username: "admin"},
        },
    })
    if err != nil {
        t.Fatal(engine)
    }

    t.Log(engine)
}

func TestGetRoleDetail(t *testing.T) {
    u, err := engine.GetRoleDetail(ctx, model.NewDescribeRequestByName("admin"))
    if err != nil {
        t.Fatal(err)
    }
    t.Log(u)
}

func TestGetRoleList(t *testing.T) {
    data, err := engine.GetRoleList(ctx, &model.RoleListOptions{ListOptions: &model.ListOptions{All: true}})
    if err != nil {
        t.Fatal(err)
    }
    t.Log(data)
}
