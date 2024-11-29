package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	. "go-security/security"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	. "go-security/erp/internal/repository"
	. "go-security/security/repository"
	. "go-security/security/service"
)

type ProfileService struct {
	UserService       *UserService
	ProfileRepository IProfileRepository
	FileUploadService IFileUploadService
}

type URL string

func NewProfileService(profileRepository IProfileRepository, fileUploadService IFileUploadService) *ProfileService {
	return &ProfileService{
		ProfileRepository: profileRepository,
		FileUploadService: fileUploadService,
	}
}

func (service *ProfileService) GetAllNotificationTypes() []NotificationType {
	return []NotificationType{NotificationTypeEmail, NotificationTypeSMS, NotificationTypeLineMessage}
}

func (service *ProfileService) GetAllProfiles(ctx context.Context) ([]*UserProfile, error) {
	return service.ProfileRepository.FindAllProfiles(ctx)
}

func (service *ProfileService) GetProfileByID(ctx context.Context, profileID uint) (*UserProfile, error) {
	profile, err := service.ProfileRepository.FindProfileByID(ctx, profileID)
	if err != nil {
		return nil, internal.ProfileNotFound
	}
	return profile, nil
}

func (service *ProfileService) GetProfileByUserID(ctx context.Context, userID uint) (*UserProfile, error) {
	profile, err := service.ProfileRepository.FindProfileByUserID(ctx, userID)
	if err != nil {
		return nil, internal.ProfileNotFound
	}
	return profile, nil
}

func (service *ProfileService) CreateDefaultProfile(ctx context.Context, userID uint) (*UserProfile, error) {
	profile := UserProfile{
		UserID: userID,
	}
	err := service.ProfileRepository.AddProfile(ctx, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (service *ProfileService) FindProfileByUserId(ctx context.Context, customerId uint) (*UserProfile, error) {
	profile, err := service.ProfileRepository.FindProfileByUserID(ctx, customerId)
	if err != nil {
		return nil, internal.ProfileNotFound
	}
	return profile, nil
}

func (service *ProfileService) AddProfile(
	ctx context.Context,
	userID uint,
	phoneNumber string,
) error {
	newProfile := UserProfile{
		UserID:      userID,
		PhoneNumber: phoneNumber,
	}
	user, err := service.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return UserNotFound
	}
	if err := service.validateProfile(&newProfile, user); err != nil {
		return err
	}
	return service.ProfileRepository.AddProfile(ctx, &newProfile)
}

func (service *ProfileService) validateProfile(profile *UserProfile, user *User) error {
	notificationTypes := []NotificationType{}
	for _, approach := range service.GetAllNotificationTypes() {
		notificationTypes = append(notificationTypes, approach)
	}

	approachSets := SetFromSlice(notificationTypes)
	if approachSets.Contains(NotificationTypeEmail) && !user.IsVerified {
		return internal.UserNotVerified
	}
	if approachSets.Contains(NotificationTypeSMS) {
		if len(profile.PhoneNumber) == 0 {
			return internal.ProfilePhoneNumberRequired
		}
		if !profile.IsPhoneNumberVerified {
			return internal.ProfilePhoneNumberNotVerified
		}
	}
	if approachSets.Contains(NotificationTypeLineMessage) && user.Platform.Name != string(PlatformLine) {
		return internal.UserPlatformNotLinePlatform
	}

	return nil
}

func (service *ProfileService) UpdateProfile(ctx context.Context, profile *UserProfile, values any) error {
	return service.ProfileRepository.UpdateProfile(ctx, profile, values)
}

// TODO: if profile is not exists, prompt user to finish profiling...
func (service *ProfileService) IsProfileExists(ctx context.Context, customerId uint) bool {
	_, err := service.ProfileRepository.FindProfileByUserID(ctx, customerId)
	return err == nil
}

func (service *ProfileService) isImageValid(file *os.File) bool {
	ext := strings.ToLower(filepath.Ext(file.Name()))
	validTypes := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff",
	}
	return slices.Contains(validTypes, ext)

}

func (service *ProfileService) UploadUserProfilePicture(ctx context.Context, userID uint, file *os.File) (*URL, *minio.UploadInfo, error) {
	if !service.isImageValid(file) {
		return nil, nil, internal.ProfileImageNotValid
	}

	profile, err := service.FindProfileByUserId(ctx, userID)
	if err != nil && profile == nil {
		log.Info().Msg("Profile not found, creating new profile")
		newProfile, err := service.CreateDefaultProfile(ctx, userID)
		if err != nil {
			return nil, nil, err
		}
		profile = newProfile
	}
	if profile == nil {
		return nil, nil, internal.ProfileNotCreated
	}

	if len(profile.ProfilePictureURL) == 0 {
		if err := service.FileUploadService.DeleteFile(ctx, profile.ProfilePictureURL); err != nil {
			log.Warn().Err(err).Msg("Failed to delete old profile image")
		}
	}

	uploadInfo, err := service.FileUploadService.UploadFile(ctx, file)
	if err != nil {
		return nil, nil, err
	}

	profile.ProfilePictureURL = uploadInfo.Key
	if err := service.UpdateProfile(ctx, profile, map[string]any{"profile_picture_url": uploadInfo.Key}); err != nil {
		return nil, nil, err
	}

	stringUrl, err := service.FileUploadService.GetFileExpiresIn(ctx, uploadInfo.Key, 5*time.Minute)
	if err != nil {
		return nil, nil, err
	}
	url := URL(stringUrl)
	return &url, uploadInfo, nil
}
