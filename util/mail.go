package util

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"net/smtp"
	"text/template"

	"gopkg.in/gomail.v2"
	"gs.lee.was/configs"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendEmail(receiver string, title string, bodyFile string, message interface{}, attachments []string) error {
	// // Sender data.
	// from := configs.GetConfigData().EmailInfo.Address
	// password := configs.GetConfigData().EmailInfo.Password

	// // smtp server configuration.
	// smtpHost := configs.GetConfigData().EmailInfo.Host
	// smtpPort := configs.GetConfigData().EmailInfo.Port

	// if attachments == nil {
	// 	attachments = make(map[string][]byte)
	// }
	// conn, err := net.Dial("tcp", smtpHost+":"+strconv.Itoa(smtpPort))
	// if err != nil {
	// 	configs.GetConfigLog().Error(err)
	// 	return err
	// }

	// c, err := smtp.NewClient(conn, smtpHost)
	// if err != nil {
	// 	configs.GetConfigLog().Error(err)
	// 	return err
	// }

	// tlsconfig := &tls.Config{
	// 	ServerName: smtpHost,
	// }

	// if err = c.StartTLS(tlsconfig); err != nil {
	// 	configs.GetConfigLog().Error(err)
	// 	return err
	// }

	// auth := LoginAuth(from, password)

	// if err = c.Auth(auth); err != nil {
	// 	configs.GetConfigLog().Error(err)
	// 	return err
	// }

	t, err := template.ParseFiles("static/" + bodyFile)
	if err != nil {
		configs.GetConfigLog().Error(err)
		return err
	}
	body := bytes.NewBuffer(nil)
	// writer := multipart.NewWriter(body)
	// boundary := writer.Boundary()
	// mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// if len(attachments) > 0 {
	// 	mimeHeaders = fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=%s\n--%s\n", boundary, boundary)
	// }
	// body.WriteString(fmt.Sprintf("Subject: %s\nTo: %s\n%s\n\n", title, strings.Join(receiver, ","), mimeHeaders))

	t.Execute(body, message)

	// if len(attachments) > 0 {
	// 	for k, v := range attachments {
	// 		body.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
	// 		body.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
	// 		body.WriteString("Content-Transfer-Encoding: base64\n")
	// 		body.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

	// 		b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
	// 		base64.StdEncoding.Encode(b, v)
	// 		body.Write(b)
	// 		body.WriteString(fmt.Sprintf("\n--%s", boundary))
	// 	}

	// 	body.WriteString("--")
	// }

	// err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, from, receiver, body.Bytes())
	// if err != nil {
	// 	configs.GetConfigLog().Error(err)
	// 	return err
	// }

	//d := gomail.NewDialer(configs.GetConfigData().EmailInfo.Host, configs.GetConfigData().EmailInfo.Port, configs.GetConfigData().EmailInfo.Address, configs.GetConfigData().EmailInfo.Password)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", configs.GetConfigData().EmailInfo.Address)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body.String())
	println(configs.GetConfigData().EmailInfo.Address)
	for _, v := range attachments {
		m.Attach(v)
	}

	println(receiver)

	var emailRaw bytes.Buffer
	m.WriteTo(&emailRaw)

	msg := ses.RawMessage{
		Data: emailRaw.Bytes(),
	}

	configs.GetConfigLog().Debug("aws Id ::  " + configs.GetConfigData().EmailInfo.Id)
	configs.GetConfigLog().Debug("aws Pw ::  " + configs.GetConfigData().EmailInfo.Password)

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region:      aws.String("ap-northeast-2"),
			Credentials: credentials.NewStaticCredentials("AKIA4QB7EBNIEPXNU6MX", "PWIZ+M0Lstw9a/wFvLOxMf6SIxW3rwMNC6Uwm7W6", ""),
		},
	})

	toAddresses := []*string{
		aws.String(receiver),
	}

	sesClient := ses.New(sess)
	_, err = sesClient.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage:   &msg,
		Destinations: toAddresses,
		Source:       aws.String(configs.GetConfigData().EmailInfo.Address),
	})

	return err

	//if err := d.DialAndSend(m); err != nil {
	//	configs.GetConfigLog().Debug(err)
	//	return err
	//}

	return nil
}
