package repository

type OrderState string

const (
	OrderStatePending  OrderState = "Pending"
	OrderStateApproved OrderState = "Approved"

	OrderStateCanceled          OrderState = "Canceled"
	OrderStateWaitingForPayment OrderState = "WaitingForPayment"
	OrderStatePaid              OrderState = "Paid"
	OrderStateShipped           OrderState = "Shipped"
)

type NotificationType string

const (
	NotificationTypeEmail       NotificationType = "Email"
	NotificationTypeSMS         NotificationType = "SMS"
	NotificationTypeLineMessage NotificationType = "LINE"
)
