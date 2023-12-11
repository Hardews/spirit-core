/**
 * @Author: Hardews
 * @Date: 2023/4/15 14:53
 * @Description:
**/

package tool

import (
	"errors"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

var ErrOfNotARealEmail = errors.New("不是一个正确的邮箱")

const (
	Sender        = "gmt-team <twentyue@qq.com>" // 发送者邮箱
	Theme         = "gmt-team Email"             // 发送的主题
	FutureTheme   = "一封过去的信[{.theme}]"
	ReplyHtml     = "<!DOCTYPE html>\n<html>\n<head>\n\t<meta charset=\"UTF-8\">\n\t<style type=\"text/css\">\n\t\tbody {\n\t\t\tfont-family: Arial, sans-serif;\n\t\t\tfont-size: 16px;\n\t\t\tline-height: 1.5;\n\t\t\tcolor: #333333;\n\t\t\tbackground-color: #f4f4f4;\n\t\t\tpadding: 0;\n\t\t\tmargin: 0;\n\t\t}\n\t\t.container {\n\t\t\twidth: 100%;\n\t\t\tmax-width: 600px;\n\t\t\tmargin: 0 auto;\n\t\t\tbackground-color: #ffffff;\n\t\t\tpadding: 20px;\n\t\t\tbox-shadow: 0 0 20px rgba(0,0,0,0.1);\n\t\t}\n\t\th1, h2, h3, h4, h5, h6 {\n\t\t\tmargin: 0 0 10px 0;\n\t\t\tcolor: #333333;\n\t\t\tline-height: 1.2;\n\t\t\tfont-weight: bold;\n\t\t}\n\t\tp {\n\t\t\tmargin: 0 0 20px 0;\n\t\t}\n\t\ta {\n\t\t\tcolor: #007bff;\n\t\t\ttext-decoration: none;\n\t\t}\n\t\ta:hover {\n\t\t\tcolor: #0056b3;\n\t\t\ttext-decoration: underline;\n\t\t}\n\t</style>\n</head>\n<body>\n\t<div class=\"container\">\n\t\t<h1>gmt-team Reply</h1>\n\t\t<p>亲爱的用户,</p>\n\t\t<p>非常感谢您的来信。我们收到了您的邮件并阅读了您的内容。</p>\n\t\t<p>{.reply}</p>\n\t\t<p>如果您有任何其他问题或需要进一步的帮助，请随时联系我们。</p>\n\t\t<p>再次感谢您的来信。</p>\n\t\t<p><strong>年年岁岁身长，负岁年年春草长。</strong></p>\n\t\t<p><strong>[gmt-team]</strong></p>\n\t</div>\n</body>\n</html>"
	CodeHtml      = "<!DOCTYPE html>\n<html lang=\"en\">\n\n<head>\n  <meta charset=\"UTF-8\">\n  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n  <title>验证邮箱</title>\n\n  <style>\n    body,html,div,ul,li,button,p,img,h1,h2,h3,h4,h5,h6 {\n      margin: 0;\n      padding: 0;\n    }\n\n    body,html {\n      background: #fff;\n      line-height: 1.8;\n    }\n\n    h1,h2,h3,h4,h5,h6 {\n      line-height: 1.8;\n    }\n\n    .email_warp {\n      height: 100vh;\n      min-height: 500px;\n      font-size: 14px;\n      color: #212121;\n      display: flex;\n      /* align-items: center; */\n      justify-content: center;\n    }\n\n    .logo {\n      margin: 3em auto;\n      width: 200px;\n      height: 60px;\n    }\n\n    h1.email-title {\n      font-size: 26px;\n      font-weight: 500;\n      margin-bottom: 15px;\n      color: #252525;\n    }\n\n    .links_btn {\n      border: 0;\n      background: #4C84FF;\n      color: #fff;\n      width: 100%;\n      height: 50px;\n      line-height: 50px;\n      font-size: 16px;\n      margin: 40px auto;\n      box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.15);\n      border-radius: 4px;\n      outline: none;\n      cursor: pointer;\n      transition: all 0.3s;\n      text-align: center;\n      display: block;\n      text-decoration: none;\n    }\n\n    .warm_tips {\n      color: #757575;\n      background: #f7f7f7;\n      padding: 20px;\n    }\n\n    .warm_tips .desc {\n      margin-bottom: 20px;\n    }\n\n    .qr_warp {\n      max-width: 140px;\n      margin: 20px auto;\n    }\n\n    .qr_warp img {\n      max-width: 100%;\n      max-height: 100%;\n    }\n\n    .email-footer {\n      margin-top: 2em;\n    }\n\n    #reset-password-email {\n      max-width: 500px;\n    }\n    #reset-password-email .accout_email {\n      color: #4C84FF;\n      display: block;\n      margin-bottom: 20px;\n    }\n  </style>\n</head>\n\n<body>\n  <section class=\"email_warp\">\n    <div id=\"reset-password-email\">\n      <div class=\"logo\">\n        <img src=\"{.logo}\" alt=\"logo\">\n      </div>\n\n      <h1 class=\"email-title\">\n        尊敬的用户您好：\n      </h1>\n      <p>您正在gmt team下的网站验证如下邮箱：</p>\n      <b class=\"accout_email\">{.address}</b>\n\n      <p>请注意，如果这不是您本人的操作，请忽略并关闭此邮件。</p>\n      <p>您的验证码为:</p>\n\n      <a class=\"links_btn\">{.code} （五分钟内有效）</a>\n\n      <div class=\"warm_tips\">\n\n        <p>如有任何疑问或无法验证，请通过如下方式与我们联系：</p>\n        <p>邮箱:twentyue@qq.com</p>\n        <p>本邮件由系统自动发送，请勿回复。</p>\n      </div>\n    </div>\n  </section>\n</body>\n</html>"
	FutureContent = "<!DOCTYPE html>\n<html lang=\"en\">\n\n<head>\n    <meta charset=\"UTF-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>Document</title>\n    <style>\n        * {\n            margin: 0;\n            padding: 0;\n        }\n\n        body,\n        html {\n            background: #fff;\n            line-height: 1.8;\n        }\n\n        h1,\n        h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            line-height: 1.8;\n        }\n\n        .content {\n            text-indent: 2em;\n        }\n\n        h1.email-title {\n            font-size: 26px;\n            font-weight: 500;\n            margin-bottom: 15px;\n            color: #252525;\n        }\n\n        .warm_tips {\n            color: #757575;\n            background: #f7f7f7;\n            padding: 20px;\n        }\n\n        .container {\n            width: 800px;\n            margin: 0 auto;\n        }\n    </style>\n</head>\n\n<body>\n    <div class=\"container\">\n        <h1 class=\"email-title\">\n            亲爱的<span>{.name}</span>：\n        </h1>\n        <h3>这封信从<span>{.addTime}</span> 送至<span>{.sendTime}</span></h3>\n        <br>\n        <b class=\"accout_email\"></b>\n\n        <p class=\"content\"></p>\n        <div class=\"warm_tips\">\n            <p>{.content}</p>\n            \n        </div>\n        <br>\n        <p>如果你是收到了其他人给你写的信，祝福你，曾有人记着你，想告诉你很多事。<br>\n            GMT-team 只是负责寄信，如果你有些话想要再告诉他/她，不妨当面去说吧 :-)</p>\n        <p>来自：<a herf=\"121321\">GMT-team</a><br>自动发送，请勿回复</p>\n    </div>\n    </div>\n</body>\n\n</html>"
)

const (
	Reply = iota + 1
	Code
	Future
)

var (
	sendEmailName     = os.Getenv("spirit_core_email_send_username")
	sendEmailPassword = os.Getenv("spirit_core_email_send_password")
)

func SendEmail(theme, address, content string, choice int) error {
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱
	em.From = Sender

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{address}

	// 设置主题
	em.Subject = theme

	switch choice {
	case Reply:
		// 回信
		replyHtml := strings.ReplaceAll(ReplyHtml, "{.reply}", content)
		em.HTML = []byte(replyHtml)
	case Code:
		// 验证码
		codeHtml := strings.ReplaceAll(CodeHtml, "{.address}", address)
		codeHtml = strings.ReplaceAll(codeHtml, "{.code}", content)
		em.HTML = []byte(codeHtml)
	case Future:
		em.HTML = []byte(content)
	}

	// 设置相关的配置
	return em.Send("smtp.qq.com:25", smtp.PlainAuth("", sendEmailName, sendEmailPassword, "smtp.qq.com"))
}

func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	return result
}
