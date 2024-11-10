package model

import (
	"encoding/json"
)

type Role struct {
	ID          int64       `json:"id"           gorm:"primarykey"`
	CreatedAt   int64       `json:"created_at"   gorm:"autoCreateTime"`
	UpdatedAt   int64       `json:"updated_at"   gorm:"autoCreateTime"`
	Rolename    string      `json:"rolename"     gorm:"type:varchar(20);not null;unique"`
	Nickname    string      `json:"nickname"     gorm:"type:varchar(20);not null;unique"`
	Status      bool        `json:"status"       gorm:"type:tinyint(1)"`
	Description string      `json:"description"  gorm:"type:varchar(255);comment:描述"`
	Parents     []string    `json:"parents"      gorm:"serializer:json;comment:当前角色的父级角色,用于继承权限"`
	Users       []*User     `json:"users"        gorm:"many2many:user_roles;foreignKey:Rolename;joinForeignKey:RoleRolename;References:Username;joinReferences:UserUsername"`
	Endpoints   []*Endpoint `json:"endpoints"    gorm:"many2many:role_endpoints;foreignKey:Rolename;joinForeignKey:RoleRolename;References:Identity;joinReferences:EndpointIdentity"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) String() string {
	data, _ := json.MarshalIndent(r, "", "")
	return string(data)
}

type RoleListOptions struct {
	*ListOptions
	Keyword  string
	Rolename string
	Nickname string
}

type RoleSet struct {
	Total int64   `json:"total"`
	Items []*Role `json:"items"`
}

type CreateRole struct {
	ID          int64    `json:"id"`
	Rolename    string   `json:"rolename"   binding:"required"`
	Nickname    string   `json:"nickname"   binding:"required"`
	Endpoints   []string `json:"endpoints"`
	Description string   `json:"description"`
	Status      bool     `json:"status"`
	Parents     []string `json:"parents"`
	Users       []string `json:"users"`
}

func (c *CreateRole) GetUsers() (data []*User) {
	for _, v := range c.Users {
		data = append(data, &User{Username: v})
	}

	return
}

func (c *CreateRole) GetEndpoints() (data []*Endpoint) {
	for _, v := range c.Endpoints {
		data = append(data, &Endpoint{Identity: v})
	}

	return
}

func NewRole(in *CreateRole) *Role {
	if len(in.Parents) == 0 {
		in.Parents = nil // 去除前端显示的[]
	}
	return &Role{
		ID:          in.ID,
		Rolename:    in.Rolename,
		Nickname:    in.Nickname,
		Status:      true,
		Description: in.Description,
		Parents:     in.Parents,
		Users:       in.GetUsers(),
		Endpoints:   in.GetEndpoints(),
	}
}

func (r *Role) SetRole(in *CreateRole) *Role {
	if in.Rolename != "" {
		r.Rolename = in.Rolename
	}
	if in.Rolename != "" {
		r.Rolename = in.Rolename
	}
	if in.Description != "" {
		r.Description = in.Description
	}
	if len(in.Endpoints) != 0 {
		r.Endpoints = in.GetEndpoints()
	}
	if len(in.Parents) == 0 {
		in.Parents = nil // 去除前端显示的[]
	}
	r.Status = in.Status
	r.Parents = in.Parents
	r.Users = in.GetUsers()

	return r
}
