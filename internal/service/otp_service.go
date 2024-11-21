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

func GenerateOtpCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10)) // Generate a digit (0-9)
	}
	return code
}

type OTP struct {
	UserId         uint
	Purpose        Purpose
	Code           string
	ExpirationTime int64
}

type OtpService struct {
	OtpCache         map[uint]map[Purpose]*OTP
	OtpLock          *sync.Mutex
	OtpGeneratorFunc func() string
}

func NewOtpService(generatorFunc func() string) *OtpService {
	return &OtpService{
		OtpCache:         make(map[uint]map[Purpose]*OTP),
		OtpLock:          new(sync.Mutex),
		OtpGeneratorFunc: generatorFunc,
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
	code := service.OtpGeneratorFunc()
	otp := &OTP{
		UserId:         userId,
		Purpose:        purpose,
		Code:           code,
		ExpirationTime: time.Now().Add(time.Minute * 5).Unix(),
	}

	service.OtpLock.Lock()
	defer service.OtpLock.Unlock()

	if service.OtpCache[userId] == nil {
		service.OtpCache[userId] = make(map[Purpose]*OTP)
	}
	service.OtpCache[userId][purpose] = otp
	return otp
}

func (service *OtpService) mustGetOtp(userId uint, purpose Purpose) (*OTP, error) {
	service.OtpLock.Lock()
	defer service.OtpLock.Unlock()

	otpMap, ok := service.OtpCache[userId]
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
