package web

import "errors"

var (
	EmailRateLimitExceeded = errors.New("EmailRateLimitExceeded")
	UnableToIdentifyUser   = errors.New("UnableToIdentifyUser")
)
