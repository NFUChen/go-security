package internal

import "errors"

var (
	UserNotFound           = errors.New("UserNotFound")
	UserAlreadyExists      = errors.New("UserAlreadyExists")
	UserPasswordNotMatched = errors.New("UserPasswordNotMatched")
	UserNameNotAllowed     = errors.New("UserNameNotAllowed")
	UserEmailNotAllowed    = errors.New("UserEmailNotAllowed")
	UserPasswordNotAllowed = errors.New("UserPasswordNotAllowed")
	UserRoleNotAllowed     = errors.New("UserRoleNotAllowed")
	UserRoleNotFound       = errors.New("UserRoleNotFound")
)
