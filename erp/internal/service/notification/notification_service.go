package notification

import . "go-security/erp/internal/repository"

type INotificationService interface {
	Name() string
	SendOrderWaitingForApprovalMessage(order *Order, profile *UserProfile) error
	SendOrderApprovedMessage(order *Order, profile *UserProfile) error
}
