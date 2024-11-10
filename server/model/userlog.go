package model

import (
    "encoding/json"
)

type Userlog struct {
    ID        int64  `json:"id"               gorm:"primarykey"`
    CreatedAt int64  `json:"created_at"       gorm:"autoCreateTime"`
    UpdatedAt int64  `json:"updated_at"       gorm:"autoCreateTime"`
    Username  string `json:"username"         gorm:"type:varchar(50);not null;comment:用户名"`
    IPAddr    string `json:"ip_addr"          gorm:"type:varchar(128);not null;comment:ip地址"`
    StartAt   int64  `json:"start_at"         gorm:"type:varchar(18);not null;comment:请求开始时间"`
    Path      string `json:"path"             gorm:"type:varchar(128);not null;comment:请求路径"`
    Method    string `json:"method"           gorm:"type:varchar(50);not null;comment:请求方法"`
    Status    int64  `json:"status"           gorm:"type:int(50);not null;comment:请求状态"`
    Duration  int64  `json:"duration"         gorm:"type:int(6);not null;comment:请求耗时(ms)"`
    Browser   string `json:"browser"          gorm:"type:varchar(128);comment:客户端"`
    ClientOS  string `json:"client_os"        gorm:"type:varchar(128);comment:客户端操作系统"`
}

func (u *Userlog) String() string {
    data, _ := json.MarshalIndent(u, "", "")
    return string(data)
}

func (Userlog) TableName() string {
    return "userlogs"
}

type UserlogSet struct {
    Total int64      `json:"total"`
    Items []*Userlog `json:"items"`
}

type UserlogListOptions struct {
    *ListOptions
    Keyword string
}
