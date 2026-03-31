package service

/*
 * @Desc: email服务
 * @author: 福狼
 * @version: v1.0.0
 */

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"time"

	"github.com/GoFurry/gofurry-user/common"
	"github.com/GoFurry/gofurry-user/common/util"
	"github.com/GoFurry/gofurry-user/roof/env"
	"gopkg.in/gomail.v2"
)

// EmailSendCode 发送邮箱验证码
func EmailSendCode(email string) (code string, gfsError common.GFError) {
	// 生成6位随机验证码
	code = util.GenerateRandomCode(common.EMAIL_CODE_LENGTH)
	m := gomail.NewMessage()
	encodedName := mimeEncode("GoFurry邮件服务")
	from := encodedName + " <" + env.GetServerConfig().Email.EmailUser + ">"
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", mimeEncode("GoFurry 邮箱验证码"))

	msg := `
	<html>
	<head>
		<meta charset="UTF-8">
		<title>GoFurry 验证码</title>
		<style>
			body { font-family: "Microsoft YaHei", "Helvetica Neue", sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
			.container { background-color: #f9f9f9; border-radius: 8px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
			.logo { color: #2c3e50; font-size: 24px; font-weight: bold; margin-bottom: 20px; display: flex; align-items: center; }
			.logo span { color: #3498db; margin-right: 8px; }
			.greeting { font-size: 18px; margin-bottom: 20px; }
			.code-box { background-color: #fff; border: 1px dashed #ddd; border-radius: 4px; padding: 20px; text-align: center; margin: 25px 0; }
			.code { font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #2c3e50; margin: 0; }
			.note { color: #666; font-size: 14px; margin: 20px 0; }
			.warning { color: #e74c3c; font-size: 13px; padding: 10px; background-color: #fef0f0; border-radius: 4px; margin-top: 15px; }
			.footer { margin-top: 30px; color: #999; font-size: 12px; text-align: center; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="logo">
				<span>🐺</span> GoFurry
			</div>
			<div class="greeting">您好！</div>
			<p>感谢您使用 GoFurry 服务，您正在进行邮箱验证操作。</p>
			<div class="code-box">
				<p class="code">[ ` + code + ` ]</p>
			</div>
			<p class="note">
				• 该验证码有效期为 <strong>5分钟</strong>，请在有效期内完成验证<br>
				• 验证码仅用于本次操作，请勿向他人泄露
			</p>
			<div class="warning">
				如果您未发起此操作，请忽略本邮件，您的账号安全不会受到影响。
			</div>
			<div class="footer">
				<p>GoFurry 邮箱服务 © ` + strconv.Itoa(time.Now().Year()) + `</p>
			</div>
		</div>
	</body>
	</html>
	`

	m.SetBody("text/html; charset=UTF-8", msg)
	d := gomail.NewDialer(
		env.GetServerConfig().Email.EmailHost,
		env.GetServerConfig().Email.EmailPort,
		env.GetServerConfig().Email.EmailUser,
		env.GetServerConfig().Email.EmailPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		gfsError = common.NewServiceError("邮件发送失败..." + err.Error())
	}
	return code, gfsError
}

// IsEmailValid 校验邮箱是否合法
func IsEmailValid(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// mimeEncode 对中文进行MIME编码
func mimeEncode(s string) string {
	// 检查是否包含非ASCII字符
	hasNonASCII := false
	for _, r := range s {
		if r > 127 {
			hasNonASCII = true
			break
		}
	}
	if !hasNonASCII {
		return s // 纯英文无需编码
	}

	// 中文使用UTF-8编码后再Base64编码
	b := []byte(s)
	encoded := base64.StdEncoding.EncodeToString(b)
	return "=?UTF-8?B?" + encoded + "?="
}
