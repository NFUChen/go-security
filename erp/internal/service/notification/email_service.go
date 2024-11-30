package notification

import (
	. "go-security/erp/internal/repository"
	"go-security/security/service"
)

type EmailService struct {
	*service.SmtpService
}

func (service EmailService) Name() string {
	return "EmailService"
}

func (service EmailService) SendOrderWaitingForApprovalMessage(order *Order, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")

}

func (service EmailService) SendOrderApprovedMessage(order *Order, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")
}

func NewEmailService(smtpService *service.SmtpService) *EmailService {
	return &EmailService{SmtpService: smtpService}
}
