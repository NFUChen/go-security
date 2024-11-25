package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	. "go-security/erp/internal/repository"
	. "go-security/erp/internal/service/notification"
	. "go-security/security/repository"
	"time"
)

type OrderService struct {
	OrderRepository IOrderRepository
	ProfileService  *ProfileService
	EmailService    INotificationService
	SmsService      INotificationService
	LineService     INotificationService
}

type Notification string

const (
	NotificationWaitingForApproval Notification = "waiting_for_approval"
	NotificationApproved           Notification = "approved"
)

func NewOrderService(
	orderRepository IOrderRepository,
	profileService *ProfileService,
	emailService INotificationService,
	smsService INotificationService,
	lineService INotificationService,
) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
		ProfileService:  profileService,
		EmailService:    emailService,
		SmsService:      smsService,
		LineService:     lineService,
	}
}

func (service *OrderService) AllOrderStates() []OrderState {
	return []OrderState{
		OrderStatePending,
		OrderStateApproved,
		OrderStatePaid,
		OrderStateUnpaid,
		OrderStateShipped,
		OrderStateCanceled,
	}
}

func (service *OrderService) FindOrdersByCustomerIDAndDate(ctx context.Context, customerID uint, datetime time.Time) ([]*CustomerOrder, error) {
	return service.OrderRepository.FindOrdersByCustomerIDAndDate(ctx, customerID, datetime)
}

func (service *OrderService) PlaceOrder(ctx context.Context, customer *User) error {
	newOrder := &CustomerOrder{
		CustomerID: customer.ID,
		OrderState: OrderStatePending,
		OrderDate:  time.Now(),
	}
	if err := service.OrderRepository.CreateOrder(ctx, newOrder); err != nil {
		return err
	}

	profile, err := service.ProfileService.FindProfileByCustomerId(customer.ID)
	if err != nil {
		return err
	}

	if err := service.sendNotifications(newOrder, profile, NotificationWaitingForApproval); err != nil {
		return err
	}

	return nil
}

func (service *OrderService) ApproveOrder(ctx context.Context, orderID uint) error {
	order, err := service.OrderRepository.FindOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.OrderState != OrderStatePending {
		return internal.PendingOrderStateRequired
	}
	order.OrderState = OrderStateApproved
	if err := service.OrderRepository.UpdateOrderState(ctx, order.ID, OrderStateApproved); err != nil {
		return err
	}
	profile, err := service.ProfileService.FindProfileByCustomerId(order.CustomerID)
	if err != nil {
		return err
	}

	if err := service.sendNotifications(order, profile, NotificationApproved); err != nil {
		return err
	}

	return nil
}

func (service *OrderService) CancelOrder(ctx context.Context, orderID uint) error {
	order, err := service.OrderRepository.FindOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.OrderState = OrderStateCanceled
	if err := service.OrderRepository.UpdateOrderState(ctx, order.ID, OrderStateCanceled); err != nil {
		return err
	}
	return nil
}

func (service *OrderService) sendNotifications(order *CustomerOrder, profile *CustomerProfile, notificationType Notification) error {
	dispatch := map[Notification]func(INotificationService, *CustomerOrder, *CustomerProfile) error{
		NotificationApproved:           INotificationService.SendOrderApprovedMessage,
		NotificationWaitingForApproval: INotificationService.SendOrderWaitingForApprovalMessage,
	}

	sendFunc, ok := dispatch[notificationType]
	if !ok {
		return internal.InvalidNotificationType
	}

	notifiers := service.GetNotificationServicesByProfile(profile)
	for _, notifier := range notifiers {
		go func(notifier INotificationService) {
			if err := sendFunc(notifier, order, profile); err != nil {
				log.Warn().Err(err).Msgf("Failed to send notification message for %v with %v", notificationType, notifier.Name())
			}
		}(notifier)
	}

	return nil
}

func (service *OrderService) GetNotificationServicesByProfile(profile *CustomerProfile) []INotificationService {
	serviceMap := map[NotificationApproach]INotificationService{
		NotificationApproachEmail:       service.EmailService,
		NotificationApproachSMS:         service.SmsService,
		NotificationApproachLineMessage: service.LineService,
	}

	notifiers := []INotificationService{}
	for _, approach := range profile.NotificationApproaches {
		if notifier, ok := serviceMap[approach]; ok {
			notifiers = append(notifiers, notifier)
		}
	}
	return notifiers
}
