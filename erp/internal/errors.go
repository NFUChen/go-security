package internal

import "errors"

var (
	OrderAlreadyShipped           = errors.New("OrderAlreadyShipped")
	OrderAlreadyCancelled         = errors.New("OrderAlreadyCancelled")
	UserNotVerified               = errors.New("UserNotVerified")
	PendingOrderStateRequired     = errors.New("PendingOrderStateRequired")
	InvalidNotificationType       = errors.New("InvalidNotificationType")
	ProfilePhoneNumberRequired    = errors.New("ProfilePhoneNumberRequired")
	UserPlatformNotLinePlatform   = errors.New("UserPlatformNotLinePlatform")
	ProfilePhoneNumberNotVerified = errors.New("ProfilePhoneNumberNotVerified")
)
