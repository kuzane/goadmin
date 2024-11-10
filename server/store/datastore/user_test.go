package datastore_test

import (
    "testing"

    "goamin/server/model"
)

func TestCreateUser(t *testing.T) {
    err := engine.CreateUser(ctx, &model.User{
        Username:    "admin",
        Nickname:    "管理员",
        Password:    "123456",
        Email:       "admin@qq.com",
        Phone:       "15888888898",
        Avatar:      "xx",
        Status:      true,
        Description: "测试账户",
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log("ok")
}

func TestDeleteUser(t *testing.T) {
    err := engine.DeleteUser(ctx, &model.User{
        ID:       2,
        Username: "joker",
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log("ok")
}

func TestUpdateUser(t *testing.T) {
    err := engine.UpdateUser(ctx, &model.User{
        ID:          1,
        Username:    "joker",
        Nickname:    "小丑小",
        Password:    "",
        Email:       "joker@qq.com",
        Phone:       "",
        Avatar:      "",
        Status:      true,
        Description: "",
        Roles:       nil,
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log("ok")
}

func TestGetUserDetail(t *testing.T) {
    u, err := engine.GetUserDetail(ctx, model.NewDescribeRequestByName("joker"))
    if err != nil {
        t.Fatal(err)
    }
    t.Log(u)
}

func TestGetUserList(t *testing.T) {
    u, err := engine.GetUserList(ctx, &model.UserListOptions{
        ListOptions: &model.ListOptions{
            All:     false,
            Page:    2,
            PerPage: 1,
        },
        Keyword:  "",
        Username: "",
        Nickname: "",
        Email:    "",
        Phone:    "",
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log(u)
}
