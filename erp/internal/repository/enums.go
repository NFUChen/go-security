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

type NotificationApproach string

const (
	NotificationApproachEmail       NotificationApproach = "Email"
	NotificationApproachSMS         NotificationApproach = "SMS"
	NotificationApproachLineMessage NotificationApproach = "LINE"
)
