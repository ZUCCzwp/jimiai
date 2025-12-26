package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig 邮箱配置
type EmailConfig struct {
	SMTPHost     string // SMTP服务器地址，如：smtp.qq.com
	SMTPPort     string // SMTP端口，如：587
	FromEmail    string // 发送邮箱地址
	FromPassword string // 发送邮箱密码或授权码
	FromName     string // 发送者名称
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(config EmailConfig, toEmail string, code string) error {
	// 设置邮件内容
	subject := "吉米AI - 邮箱验证码"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您好！</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
			<p>验证码有效期为1分钟，请勿泄露给他人。</p>
			<p>如非本人操作，请忽略此邮件。</p>
			<br>
			<p>此邮件由系统自动发送，请勿回复。</p>
		</body>
		</html>
	`, code)

	// 构建邮件消息
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s<%s>\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n",
		toEmail, config.FromName, config.FromEmail, subject, body))

	// 设置SMTP认证
	auth := smtp.PlainAuth("", config.FromEmail, config.FromPassword, config.SMTPHost)

	// 发送邮件
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{toEmail}, msg)
	if err != nil {
		// 处理 "short response" 错误：这是 Go 标准库的已知问题
		// 某些 SMTP 服务器（如 QQ 邮箱）在关闭连接时可能返回不完整的响应
		// 但邮件实际上已经发送成功，可以忽略这个错误
		errStr := err.Error()
		if strings.Contains(errStr, "short response") {
			// 邮件已成功发送，只是读取响应时出错，可以忽略
			return nil
		}
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// SendUpdatePasswordEmailCode 发送修改密码邮箱验证码
func SendUpdatePasswordEmailCode(config EmailConfig, toEmail string, code string) error {
	// 设置邮件内容
	subject := "吉米AI - 修改密码验证码"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您好！</h2>
			<p>您正在修改密码，验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
			<p>验证码有效期为1分钟，请勿泄露给他人。</p>
			<p>如非本人操作，请立即修改密码并联系客服。</p>
			<br>
			<p>此邮件由系统自动发送，请勿回复。</p>
		</body>
		</html>
	`, code)

	// 构建邮件消息
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s<%s>\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n",
		toEmail, config.FromName, config.FromEmail, subject, body))

	// 设置SMTP认证
	auth := smtp.PlainAuth("", config.FromEmail, config.FromPassword, config.SMTPHost)

	// 发送邮件
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{toEmail}, msg)
	if err != nil {
		// 处理 "short response" 错误：这是 Go 标准库的已知问题
		// 某些 SMTP 服务器（如 QQ 邮箱）在关闭连接时可能返回不完整的响应
		// 但邮件实际上已经发送成功，可以忽略这个错误
		errStr := err.Error()
		if strings.Contains(errStr, "short response") {
			// 邮件已成功发送，只是读取响应时出错，可以忽略
			return nil
		}
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}
