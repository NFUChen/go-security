package security

import "errors"

var (
	UserPlatformEmpty       = errors.New("UserPlatformEmpty")
	UserNotFound            = errors.New("UserNotFound")
	UserAlreadyExists       = errors.New("UserAlreadyExists")
	UserAlreadyVerified     = errors.New("UserAlreadyVerified")
	UserPasswordNotMatched  = errors.New("UserPasswordNotMatched")
	UserNameNotAllowed      = errors.New("UserNameNotAllowed")
	UserEmailNotAllowed     = errors.New("UserEmailNotAllowed")
	UserPasswordNotAllowed  = errors.New("UserPasswordNotAllowed")
	UserRoleNotAllowed      = errors.New("UserRoleNotAllowed")
	UserRoleNotFound        = errors.New("UserRoleNotFound")
	TokenExpired            = errors.New("TokenExpired")
	TokenInvalid            = errors.New("TokenInvalid")
	OtpNotFound             = errors.New("OtpNotFound")
	OtpIncorrect            = errors.New("OtpIncorrect")
	OtpExpired              = errors.New("OtpExpired")
	ResetPasswordNotMatched = errors.New("ResetPasswordNotMatched")
)
