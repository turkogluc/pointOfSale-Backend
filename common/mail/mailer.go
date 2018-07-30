package mail

import (
	"net/smtp"
	"bytes"
	"html/template"
	"sync"
	"fmt"
	"github.com/scorredoira/email"
	"net/mail"
	. "iugo.fleet/common/logger"
	"strings"
	"time"
)

const templateFilePathSuffix = ".html"


var instance *CustomMail
var once sync.Once

func GetInstance() *CustomMail {
	once.Do(func() {
		instance = &CustomMail{}
	})
	return instance
}

type Credentials struct {
	auth		smtp.Auth
	smtpAddress	string
	smtpPort	string
	user		string
	password	string
}

type CustomMail struct {
	credentials 	Credentials
	templateFilePathPrefix string
	To      	[]string
	Bcc     	[]string
	Subject 	string
	body    	string
	attachPath      string
}

func (customMail *CustomMail) Init(user string, password string, templateFilePathPrefix, smtpAddress string, smtpPort string) {
	LogInfo("Initializing custom e-mailer.")
	customMail.credentials = Credentials{
		auth: smtp.PlainAuth("", user, password, smtpAddress),
		user: user,
		password: password,
		smtpAddress: smtpAddress,
		smtpPort: smtpPort,
	}
	customMail.templateFilePathPrefix = templateFilePathPrefix
}

func (customMail *CustomMail) SendEmailWithTemplate(to []string, bcc []string, subject string, templateFileName string, templateData interface{}) error {
	customMail.To = to
	customMail.Bcc = bcc
	customMail.Subject = subject

	if err := customMail.parseTemplate(templateFileName, templateData); err != nil {
		return err
	}

	//err := customMail.sendEmail()
	err := customMail.tryUntilSuccess(120, 10)
	if err != nil {
		return err
	}

	return nil
}

func (customMail *CustomMail) SendEmailWithTemplateAndAttachment(to []string, subject string, templateFileName string, templateData interface{}, attachPath string) error {
	customMail.To = to
	customMail.Subject = subject

	fmt.Println(" Attachment to:",to)
	if err := customMail.parseTemplate(templateFileName, templateData); err != nil {
		return err
	}

	customMail.attachPath = attachPath

	//err := customMail.sendEmail()
	err := customMail.tryUntilSuccess(120, 10)
	if err != nil {
		return err
	}

	return nil
}

func (customMail *CustomMail) SendSimpleEmail(to []string, subject string, text string) error {
	customMail.To = to
	customMail.Subject = subject

	customMail.body = text

	err := customMail.sendEmail()
	if err != nil {
		return err
	}

	return nil

}

func (customMail *CustomMail) sendEmail() error {
	LogInfo("Sending e-mail.")
	var address bytes.Buffer

	m := email.NewHTMLMessage(customMail.Subject, customMail.body)
	m.From = mail.Address{Name: "iUGO Teknoloji", Address: customMail.credentials.user}
	m.To = customMail.To
	m.Bcc = customMail.Bcc
	if strings.Compare(customMail.attachPath, "") != 0{
		m.Attach(customMail.attachPath)
	}
	address.WriteString(customMail.credentials.smtpAddress)
	address.WriteString(":")
	address.WriteString(customMail.credentials.smtpPort)

	auth := smtp.PlainAuth("", customMail.credentials.user, customMail.credentials.password, customMail.credentials.smtpAddress)

	if err := email.Send(address.String(), auth, m); err != nil {
		return err
	}
	return nil

}

func (mail *CustomMail) parseTemplate(templateFileName string, templateData interface{}) error {
	LogInfo("Parsing mail template.")
	message := new(bytes.Buffer)
	var templateFilePath bytes.Buffer
	templateFilePath.WriteString(mail.templateFilePathPrefix)
	templateFilePath.WriteString(templateFileName)
	templateFilePath.WriteString(templateFilePathSuffix)

	funcMapDivide := template.FuncMap{"minuteDivide": func(a int) string {
		return fmt.Sprintf("%02d", a / 60)
	}}

	funcMapMod := template.FuncMap{"secondMod": func(a int) string {
		return fmt.Sprintf("%02d", a % 60)
	}}

	funcMapMinuteFormat := template.FuncMap{
		"minuteFormat": func(i int) string { return fmt.Sprintf("%02d", i) },
	}

	funcMapSum := template.FuncMap{"sumByOne": func(a int) int {
		return (a + 1)
	}}

	funcMapTruncateFloat := template.FuncMap{"truncateFloat": func(a float32) string {
		return fmt.Sprintf("%.1f", a)
	}}

	funcMapFormatDate := template.FuncMap{"formatDate": func(a int64) string {
		a = a/1000
		tm := time.Unix(a, 0)
		return tm.Format("02/01/2006 15:04:05")

	}}

	funcMapFormatDateOnly := template.FuncMap{"formatDateOnly": func(tm time.Time) string {
		return tm.Format("02/01/2006")

	}}

	funcMapFormatValue := template.FuncMap{"formatValue": func(a float64) string{
		if a > 10{
			return fmt.Sprintf("%.2f km/h", a)
		}
		return fmt.Sprintf("%.2f g", a)
	}}

	funcMapFormatDuration := template.FuncMap{"formatDuration": func(a int64) string{
		return fmt.Sprintf("%.02d s", a)
	}}

	funcMapFormatLimit := template.FuncMap{"formatLimit": func(a float64) string{
		if a == 0 {
			return "-"
		}

		return fmt.Sprintf("%.2f km/h", a)

	}}

	t := template.Must(template.New(templateFileName + ".html").
		Funcs(funcMapFormatDate).
		Funcs(funcMapFormatDateOnly).
		Funcs(funcMapTruncateFloat).
		Funcs(funcMapSum).
		Funcs(funcMapMinuteFormat).
		Funcs(funcMapMod).
		Funcs(funcMapDivide).
		Funcs(funcMapFormatDuration).
		Funcs(funcMapFormatValue).
		Funcs(funcMapFormatLimit).
		ParseFiles(templateFilePath.String()))

	if err := t.Execute(message, templateData); err != nil {
		return err
	}


	mail.body = message.String()
	return nil
}

func (customMail *CustomMail)tryUntilSuccess(waitSec time.Duration, repeatCount int) (err error) {
	err = nil
	counter := 0
	for {
		counter++
		if repeatCount > 0 {
			if counter > repeatCount {
				break
			}
		}
		err = customMail.sendEmail()
		if err != nil {
			LogError(err)
			LogInfo("Email could not be sent. Waiting for", waitSec, "seconds to retry...")
			time.Sleep(waitSec * time.Second)
			continue
		}
		break
	}
	return
}
