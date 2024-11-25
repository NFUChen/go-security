package notification

import . "go-security/erp/internal/repository"

type INotificationService interface {
	Name() string
	SendOrderWaitingForApprovalMessage(order *CustomerOrder, profile *CustomerProfile) error
	SendOrderApprovedMessage(order *CustomerOrder, profile *CustomerProfile) error
}
