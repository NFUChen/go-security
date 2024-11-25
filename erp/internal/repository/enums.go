package repository

type OrderState string

const (
	OrderStatePending  OrderState = "Pending"
	OrderStateApproved OrderState = "Approved"

	OrderStateCanceled OrderState = "Canceled"

	OrderStatePaid    OrderState = "Paid"
	OrderStateUnpaid  OrderState = "Unpaid"
	OrderStateShipped OrderState = "Shipped"
)

type NotificationApproach string

const (
	NotificationApproachEmail       NotificationApproach = "Email"
	NotificationApproachSMS         NotificationApproach = "SMS"
	NotificationApproachLineMessage NotificationApproach = "LINE"
)
