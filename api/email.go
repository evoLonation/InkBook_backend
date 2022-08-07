package api

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
)

func SendEmail(address []string, vCode string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("InkBook <413935740@qq.com>")
	e.To = address
	t := time.Now().Format("2006-01-02 15:04:05")
	//设置文件发送的内容
	content := fmt.Sprintf(`<div>
			<div>
				<img src="https://img1.baidu.com/it/u=2394665572,2612107105&fm=253&fmt=auto&app=138&f=JPEG?w=100&h=100" width=100px height=100px>
			</div>
			<div>
				尊敬的%s，您好！
			</div>
			<div style="padding: 8px 40px 8px 50px;">
				<p>您于 %s 提交的邮箱验证，本次验证码为<u><strong>%s</strong></u>。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
			</div>
			<div>
				<p>此邮箱为系统邮箱，请勿回复。</p>
			</div>
		</div>`, address[0], t, vCode)
	e.Text = []byte(content)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "413935740@qq.com", "ukdwwhkaegvpcbch", "smtp.qq.com"))
	return err
}
