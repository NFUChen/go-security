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
	ProfileNotFound               = errors.New("ProfileNotFound")
	ProfileImageNotValid          = errors.New("ProfileImageNotValid")
	UnableToConvertFile           = errors.New("UnableToConvertFile")
	ProfileNotCreated             = errors.New("ProfileNotCreated")
	ProfileAlreadyExists          = errors.New("ProfileAlreadyExists")
	DefaultPolicyNotFound         = errors.New("DefaultPolicyNotFound")
	ProductNameRequired           = errors.New("ProductNameRequired")
	ProductCategoryRequired       = errors.New("ProductCategoryRequired")
	CategoryNameRequired          = errors.New("CategoryNameRequired")
	ProductNotFound               = errors.New("ProductNotFound")
	ProductCategoryNotFound       = errors.New("ProductCategoryNotFound")
	CategoryAlreadyExists         = errors.New("CategoryAlreadyExists")
	CategoryContainsProducts      = errors.New("CategoryContainsProducts")
	PricingPolicyAlreadyExists    = errors.New("PricingPolicyAlreadyExists")
	PricingPolicyNotFound         = errors.New("PricingPolicyNotFound")
)
