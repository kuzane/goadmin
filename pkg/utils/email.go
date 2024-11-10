package utils

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
)

type EmailBodyData struct {
	Name    string // 用户Nickname
	Content string // 通知的内容的变量
	Domain  string // 域名地址
}

func SendEmail(host string, port int, from, pass string, toEmail []string, toCopy []string, subject, emailTemplate string, data EmailBodyData) error {
	tmpl, err := template.New("email").Parse(emailTemplate) // 渲染邮件模板
	if err != nil {
		return fmt.Errorf("模板渲染失败: %v", err)
	}
	var body bytes.Buffer // 创建一个邮件内容缓存
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("模板exec失败: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(from, "goAdmin邮件通知")) //添加别名，即“XX官方”
	m.SetHeader("To", toEmail...)                                 //发送用户
	m.SetHeader("Cc", toCopy...)                                  //抄送用户
	m.SetHeader("Subject", subject)                               //设置邮件主题
	m.SetBody("text/html", body.String())                         //设置邮件正文
	d := gomail.NewDialer(host, port, from, pass)
	if err := d.DialAndSend(m); err != nil { // 发送邮件
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	return err
}
