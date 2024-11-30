package notification

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	. "go-security/erp/internal/repository"
)

type AwsSnsService struct {
	Client *sns.Client
}

func (service *AwsSnsService) Name() string {
	return "AwsSnsService"
}

func (service *AwsSnsService) SendOrderWaitingForApprovalMessage(order *Order, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")
}

func (service *AwsSnsService) SendOrderApprovedMessage(order *Order, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")
}

func NewAwsSnsService(client *sns.Client) *AwsSnsService {
	service := &AwsSnsService{
		Client: client,
	}

	return service
}

func (service *AwsSnsService) PostConstruct() {}

func (service *AwsSnsService) SendShortMessage(ctx context.Context, phoneNumber string, message string) (*sns.PublishOutput, error) {
	return service.Client.Publish(ctx, &sns.PublishInput{
		Message:     &message,
		PhoneNumber: &phoneNumber,
	})

}
