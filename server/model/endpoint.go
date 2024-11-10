package model

import "encoding/json"

// api接口
type Endpoint struct {
    ID        int64   `json:"id"           gorm:"primarykey"`
    CreatedAt int64   `json:"created_at"   gorm:"autoCreateTime"`
    UpdatedAt int64   `json:"updated_at"   gorm:"autoCreateTime"`
    Path      string  `json:"path"         gorm:"type:varchar(255);index:idx_path_method,unique"`
    Method    string  `json:"method"       gorm:"type:varchar(255);index:idx_path_method,unique"`
    Module    string  `json:"module"       gorm:"type:varchar(255);comment:模块(系统管理\工单系统\堡垒机系统)"`
    Kind      string  `json:"kind"         gorm:"type:varchar(255);comment:类别(系统模块下的子类别)"`
    Identity  string  `json:"identity"     gorm:"type:varchar(64);index:idx_identity,unique;comment:接口的唯一标识,做为和Roles的外键关联"`
    Remark    string  `json:"remark"       gorm:"type:varchar(255);comment:描述"`
    Roles     []*Role `json:"roles"        gorm:"many2many:role_endpoints;foreignKey:Identity;joinForeignKey:EndpointIdentity;References:Rolename;joinReferences:RoleRolename"`
}

func (Endpoint) TableName() string {
    return "endpoint"
}

func (v *Endpoint) String() string {
    data, _ := json.MarshalIndent(v, "", "")
    return string(data)
}

type EndpointListOptions struct {
    *ListOptions
    Keyword  string
    Kind     string
    Path     string
    Method   string
    Module   string
    Identity string
    Remark   string
}

func NewEndpointListAll() *EndpointListOptions {
    return &EndpointListOptions{ListOptions: NewAllListOptions()}
}

type EndpointSet struct {
    Total int64       `json:"total"`
    Items []*Endpoint `json:"items"`
}

type CreateEndpoint struct {
    ID       int64    `json:"id"`
    Method   string   `json:"method"`
    Path     string   `json:"path"`
    Module   string   `json:"module"`
    Kind     string   `json:"kind"`
    Identity string   `json:"identity"`
    Remark   string   `json:"remark"`
    Roles    []string `json:"roles" `
}

func NewEndpoint(in *CreateEndpoint) *Endpoint {
    return &Endpoint{
        ID:       in.ID,
        Method:   in.Method,
        Path:     in.Path,
        Module:   in.Module,
        Kind:     in.Kind,
        Identity: in.Identity,
        Remark:   in.Remark,
        Roles:    in.GetRoles(),
    }
}

func (c *CreateEndpoint) GetRoles() (data []*Role) {
    for _, v := range c.Roles {
        data = append(data, &Role{Rolename: v})
    }

    return data
}
