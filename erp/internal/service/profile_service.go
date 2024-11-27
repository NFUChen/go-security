package service

import (
	"context"
	"go-security/erp/internal"
	. "go-security/erp/internal/repository"
	. "go-security/security/repository"
	. "go-security/security/service"
)

type ProfileService struct {
	UserService       *UserService
	ProfileRepository IProfileRepository
}

func NewProfileService(profileRepository IProfileRepository) *ProfileService {
	return &ProfileService{
		ProfileRepository: profileRepository,
	}
}

func (service *ProfileService) GetProfileByID(ctx context.Context, profileID uint) (*UserProfile, error) {
	return service.ProfileRepository.FindProfileByID(ctx, profileID)
}

func (service *ProfileService) GetProfileByUserID(ctx context.Context, userID uint) (*UserProfile, error) {
	return service.ProfileRepository.FindProfileByUserId(ctx, userID)
}

func (service *ProfileService) FindProfileByCustomerId(ctx context.Context, customerId uint) (*UserProfile, error) {
	return service.ProfileRepository.FindProfileByUserId(ctx, customerId)
}

func (service *ProfileService) AddProfile(
	ctx context.Context,
	userID uint,
	notificationApproaches []NotificationApproach,
	phoneNumber string,
) error {
	newProfile := UserProfile{
		UserID:                 userID,
		NotificationApproaches: notificationApproaches,
		PhoneNumber:            phoneNumber,
	}
	user, err := service.UserService.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := service.validateProfile(&newProfile, user); err != nil {
		return err
	}
	return service.ProfileRepository.AddProfile(ctx, &newProfile)
}

func (service *ProfileService) validateProfile(profile *UserProfile, user *User) error {
	approachSets := internal.SetFromSlices(profile.NotificationApproaches)
	if approachSets.Contains(NotificationApproachEmail) && !user.IsVerified {
		return internal.UserNotVerified
	}
	if approachSets.Contains(NotificationApproachSMS) {
		if len(profile.PhoneNumber) == 0 {
			return internal.ProfilePhoneNumberRequired
		}
		if !profile.IsPhoneNumberVerified {
			return internal.ProfilePhoneNumberNotVerified
		}
	}
	if approachSets.Contains(NotificationApproachLineMessage) && user.Platform.Name != string(PlatformLine) {
		return internal.UserPlatformNotLinePlatform
	}

	return nil
}

func (service *ProfileService) UpdateProfile(ctx context.Context, profile *UserProfile, values any) error {
	return service.ProfileRepository.UpdateProfile(ctx, profile, values)
}

// TODO: if profile is not exists, prompt user to finish profiling...
func (service *ProfileService) IsProfileExists(ctx context.Context, customerId uint) bool {
	_, err := service.ProfileRepository.FindProfileByUserId(ctx, customerId)
	return err == nil
}
