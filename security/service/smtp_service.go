package service

import (
	"gopkg.in/gomail.v2"
	"os"
)

type ContentType string

const (
	ContentTypeText = "text/plain"
	ContentTypeHtml = "text/html"
)

type ISmtpService interface {
	CreateNewMessage(to string, subject string, body string, contentType ContentType, attachments ...*os.File) *gomail.Message
	SendEmail(message *gomail.Message) error
	GetSmtpConfig() *SmtpConfig
}

type SmtpConfig struct {
	CompanyName    string `yaml:"company_name"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	SenderEmail    string `yaml:"sender_email"`
	SenderPassword string `yaml:"sender_password" json:"-"`
}

type SmtpService struct {
	SmtpConfig *SmtpConfig
	Dialer     *gomail.Dialer
}

func (service *SmtpService) PostConstruct() {}

func (service *SmtpService) GetSmtpConfig() *SmtpConfig {
	return service.SmtpConfig
}

func NewSmtpService(config *SmtpConfig) *SmtpService {
	dialer := gomail.NewDialer(
		config.Host,
		config.Port,
		config.SenderEmail,
		config.SenderPassword,
	)

	service := &SmtpService{
		SmtpConfig: config,
		Dialer:     dialer,
	}

	return service
}

func (service *SmtpService) CreateNewMessage(to string, subject string, body string, contentType ContentType, attachments ...*os.File) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", service.SmtpConfig.SenderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody(string(contentType), body)

	for _, attachment := range attachments {
		message.Attach(attachment.Name())
	}
	return message
}

func (service *SmtpService) SendEmail(message *gomail.Message) error {
	return service.Dialer.DialAndSend(message)
}
