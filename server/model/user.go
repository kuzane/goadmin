package model

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64   `json:"id"             gorm:"primarykey"`
	CreatedAt    int64   `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt    int64   `json:"updated_at"     gorm:"autoCreateTime"`
	Username     string  `json:"username"       gorm:"type:varchar(50);not null;unique;comment:用户名"`
	Password     string  `json:"-"              gorm:"type:varchar(255);column:password;not null;comment:用户密码"`
	Nickname     string  `json:"nickname"       gorm:"type:varchar(50);comment:中文名"`
	Email        string  `json:"email"          gorm:"type:varchar(100);not null;unique;comment:邮箱"`
	Phone        string  `json:"phone"          gorm:"type:varchar(15);not null;unique;comment:手机号"`
	Avatar       string  `json:"avatar"         gorm:"type:varchar(255);comment:头像"`
	Status       bool    `json:"status"         gorm:"type:tinyint(1);default:1;comment:状态"`
	Description  string  `json:"description"    gorm:"type:varchar(255);comment:描述"`
	AccessToken  string  `json:"-"              gorm:"type:varchar(255);column:access_token;"`
	RefreshToken string  `json:"-"              gorm:"type:varchar(255);column:refresh_token;"`
	Roles        []*Role `json:"roles"          gorm:"many2many:user_roles;foreignKey:Username;joinForeignKey:UserUsername;References:Rolename;joinReferences:RoleRolename"`
}

func (u *User) String() string {
	data, _ := json.MarshalIndent(u, "", "")
	return string(data)
}

func (User) TableName() string {
	return "users"
}

type UserSet struct {
	Total int64   `json:"total"`
	Items []*User `json:"items"`
}

type UserListOptions struct {
	*ListOptions
	Keyword  string
	Username string
	Nickname string
	Email    string
	Phone    string
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUser struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"    binding:"required"`
	Nickname    string   `json:"nickname"    binding:"required"`
	Email       string   `json:"email"       binding:"required"`
	Phone       string   `json:"phone"       binding:"required"`
	Roles       []string `json:"roles"       binding:"required"`
	Password    string   `json:"password"`
	Status      bool     `json:"status"`
	Description string   `json:"description"`
}

func (c *CreateUser) GetRoles() (data []*Role) {
	for _, rolename := range c.Roles {
		data = append(data, &Role{Rolename: rolename})
	}

	return
}

func NewUser(in *CreateUser) (*User, error) {
	u := &User{
		ID:          in.ID,
		Username:    in.Username,
		Password:    in.Password,
		Nickname:    in.Nickname,
		Email:       in.Email,
		Phone:       in.Phone,
		Avatar:      "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif", //统一使用系统默认头像
		Status:      true,
		Description: in.Description,
		Roles:       in.GetRoles(),
	}

	if err := u.PasswordEncryption(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) SetUser(in *CreateUser) (*User, error) {
	if in.Username != "" {
		u.Username = in.Username
	}
	if in.Nickname != "" {
		u.Nickname = in.Nickname
	}
	if in.Password != "" {
		u.Password = in.Password
	}

	// 对密码进行加密
	if err := u.PasswordEncryption(); err != nil {
		return nil, err
	}

	if in.Email != "" {
		u.Email = in.Email
	}
	if in.Phone != "" {
		u.Phone = in.Phone
	}
	if in.Description != "" {
		u.Description = in.Description
	}
	u.Status = in.Status
	if len(in.Roles) > 0 {
		u.Roles = in.GetRoles()
	}

	return u, nil
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"     binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (u *User) SetUserPasswordByChangePasswordRequest(in *ChangePasswordRequest) (*User, error) {
	if in.NewPassword != in.ConfirmPassword {
		return nil, fmt.Errorf("两次输入的密码不相等！")
	}

	if err := u.PasswordVerifiers(in.OldPassword); err != nil {
		return nil, fmt.Errorf("旧密码错误!")
	}

	if in.NewPassword != "" {
		u.Password = in.NewPassword
	}

	if err := u.PasswordEncryption(); err != nil {
		return nil, err
	}

	return u, nil
}

type ChangeProfileRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (u *User) SetUserPasswordByChangeProfileRequest(in *ChangeProfileRequest) *User {
	if in.Email != "" {
		u.Email = in.Email
	}
	if in.Phone != "" {
		u.Phone = in.Phone
	}

	if in.Description != "" {
		u.Description = in.Description
	}

	return u
}

type DESCRIBEUSER_BY int64

const (
	DESCRIBEUSER_BY_ID    DESCRIBEUSER_BY = 0
	DESCRIBEUSER_BY_NAME  DESCRIBEUSER_BY = 1
	DESCRIBEUSER_BY_AK    DESCRIBEUSER_BY = 2
	DESCRIBEUSER_BY_EMAIL DESCRIBEUSER_BY = 3
	DESCRIBEUSER_BY_PHONE DESCRIBEUSER_BY = 4
)

// 获取用户详情，除默认的id name之外还有其他的选项
type DetailUserRequest struct {
	DescribeBy DESCRIBEUSER_BY
	Id         string // id
	Name       string // username
	AK         string // access_token
	Email      string
	Phone      string
}

func NewDetailUserRequestByID(id string) *DetailUserRequest {
	return &DetailUserRequest{
		DescribeBy: DESCRIBEUSER_BY_ID,
		Id:         id,
	}
}

func NewNewDetailUserRequestByName(name string) *DetailUserRequest {
	return &DetailUserRequest{
		DescribeBy: DESCRIBEUSER_BY_NAME,
		Name:       name,
	}
}

func NewNewDetailUserRequestByAK(ak string) *DetailUserRequest {
	return &DetailUserRequest{
		DescribeBy: DESCRIBEUSER_BY_AK,
		AK:         ak,
	}
}

func NewNewDetailUserRequestByEmail(mail string) *DetailUserRequest {
	return &DetailUserRequest{
		DescribeBy: DESCRIBEUSER_BY_EMAIL,
		Email:      mail,
	}
}

func NewNewDetailUserRequestByPhone(phone string) *DetailUserRequest {
	return &DetailUserRequest{
		DescribeBy: DESCRIBEUSER_BY_PHONE,
		Phone:      phone,
	}
}

// PasswordEncryption 密码加密
func (u *User) PasswordEncryption() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {

		return err
	}
	u.Password = string(hash)

	return nil
}

// PasswordVerifiers 密码加密
func (u *User) PasswordVerifiers(pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		return err
	}

	return nil
}
