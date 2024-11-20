package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"os"
	"sync"
	"time"
)

type ContentType string

const (
	ContentTypeText = "text/plain"
	ContentTypeHtml = "text/html"
)

type ISmtpService interface {
	CreateNewMessage(to string, subject string, body string, contentType ContentType, attachments ...*os.File) *gomail.Message
	SendEmail(message *gomail.Message)
	GetSmtpConfig() *SmtpConfig
}

type SmtpConfig struct {
	CompanyName    string `yaml:"company_name"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Sender         string `yaml:"sender"`
	SenderPassword string `yaml:"sender_password"`
}

type SmtpService struct {
	SmtpConfig          *SmtpConfig
	Dialer              *gomail.Dialer
	MessageQueue        []*gomail.Message
	MessageChannel      chan *gomail.Message
	MessageQueueContext context.Context
	QueueLock           *sync.Mutex
}

func (service *SmtpService) GetSmtpConfig() *SmtpConfig {
	return service.SmtpConfig
}

func NewSmtpService(ctx context.Context, config *SmtpConfig) *SmtpService {
	dialer := gomail.NewDialer(
		config.Host,
		config.Port,
		config.Sender,
		config.SenderPassword,
	)

	service := &SmtpService{
		SmtpConfig:          config,
		Dialer:              dialer,
		MessageQueue:        []*gomail.Message{},
		MessageChannel:      make(chan *gomail.Message),
		QueueLock:           new(sync.Mutex),
		MessageQueueContext: ctx,
	}

	go service.keepPopMessageQueue()
	go service.listenAndHandleEmailQueue()
	return service
}

func (service *SmtpService) CreateNewMessage(to string, subject string, body string, contentType ContentType, attachments ...*os.File) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", service.SmtpConfig.Sender)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody(string(contentType), body)

	for _, attachment := range attachments {
		message.Attach(attachment.Name())
	}
	return message
}

func (service *SmtpService) SendEmail(message *gomail.Message) {
	service.QueueLock.Lock()
	defer service.QueueLock.Unlock()
	service.MessageQueue = append(service.MessageQueue, message)
}

func (service *SmtpService) keepPopMessageQueue() {
	log.Info().Msg("Starting to pop message queue")
	for {
		select {
		case <-service.MessageQueueContext.Done():
			return
		default:
			time.Sleep(1 * time.Second)
			if len(service.MessageQueue) == 0 {
				continue
			}
			log.Info().Msgf("Sending email message: %v", service.MessageQueue[0])
			service.MessageChannel <- service.MessageQueue[0]
			service.QueueLock.Lock()
			service.MessageQueue = service.MessageQueue[1:]
			service.QueueLock.Unlock()
		}
	}
}

func (service *SmtpService) listenAndHandleEmailQueue() {
	log.Info().Msg("Starting to listen and handle email queue")
	for {
		select {
		case <-service.MessageQueueContext.Done():
			return
		case message := <-service.MessageChannel:
			log.Info().Msgf("Sending email message: %v", message)
			if err := service.Dialer.DialAndSend(message); err != nil {
				log.Warn().Msgf("Failed to send email: %v", err)
				service.QueueLock.Lock()
				log.Info().Msgf("Re-queue email message: %v", message)
				service.MessageQueue = append(service.MessageQueue, message)
				service.QueueLock.Unlock()
			}
		}
	}
}
