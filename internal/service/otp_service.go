package service

import (
	"fmt"
	"go-security/internal"
	"math/rand"
	"sync"
	"time"
)

type Purpose string

const (
	PurposeGuestEmailVerification Purpose = "guest_email_verification"
	PurposeResetPassword          Purpose = "reset_password"
)

type OTP struct {
	UserId         uint
	Purpose        Purpose
	Code           string
	ExpirationTime int64
}

type OtpService struct {
	otpCache map[uint]map[Purpose]*OTP
	otpLock  *sync.Mutex
}

func NewOtpService() *OtpService {
	return &OtpService{
		otpCache: make(map[uint]map[Purpose]*OTP),
		otpLock:  new(sync.Mutex),
	}
}

func (service *OtpService) generateOtp() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10)) // Generate a digit (0-9)
	}
	return code
}

func (service *OtpService) GenerateOtp(userId uint, purpose Purpose) *OTP {
	code := service.generateOtp()
	otp := &OTP{
		UserId:         userId,
		Purpose:        purpose,
		Code:           code,
		ExpirationTime: time.Now().Add(time.Minute * 5).Unix(),
	}

	service.otpLock.Lock()
	defer service.otpLock.Unlock()

	if service.otpCache[userId] == nil {
		service.otpCache[userId] = make(map[Purpose]*OTP)
	}
	service.otpCache[userId][purpose] = otp
	return otp
}

func (service *OtpService) mustGetOtp(userId uint, purpose Purpose) (*OTP, error) {
	service.otpLock.Lock()
	defer service.otpLock.Unlock()

	otpMap, ok := service.otpCache[userId]
	if !ok {
		return nil, internal.OtpNotFound
	}

	otp, ok := otpMap[purpose]
	if !ok {
		return nil, internal.OtpNotFound
	}

	return otp, nil
}

func (service *OtpService) GetOtp(userId uint, purpose Purpose) (*OTP, error) {
	return service.mustGetOtp(userId, purpose)
}

func (service *OtpService) VerifyOtp(userId uint, purpose Purpose, code string) error {
	cachedOtp, err := service.mustGetOtp(userId, purpose)
	if err != nil {
		return err
	}

	if cachedOtp.Code != code {
		return internal.OtpIncorrect
	}
	return nil
}
