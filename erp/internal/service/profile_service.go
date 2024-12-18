package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	. "go-security/security"
	"gorm.io/gorm"
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
	Engine                      *gorm.DB
	PricingPolicyService        *PricingPolicyService
	UserService                 *UserService
	NotificationApproachService *NotificationApproachService
	ProfileRepository           IProfileRepository
	FileUploadService           IFileUploadService
}

func NewProfileService(
	Engine *gorm.DB,
	userService *UserService,
	profileRepository IProfileRepository,
	fileUploadService IFileUploadService,
	notificationApproachService *NotificationApproachService,
) *ProfileService {
	return &ProfileService{
		Engine:                      Engine,
		UserService:                 userService,
		ProfileRepository:           profileRepository,
		FileUploadService:           fileUploadService,
		NotificationApproachService: notificationApproachService,
	}
}

func (service *ProfileService) InjectPricingPolicyService(pricingPolicyService *PricingPolicyService) {
	service.PricingPolicyService = pricingPolicyService
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
	policy, err := service.PricingPolicyService.GetDefaultPolicy(ctx)
	if err != nil {
		return nil, internal.DefaultPolicyNotFound
	}

	profile := UserProfile{
		UserID:          userID,
		PricingPolicyID: policy.ID,
	}
	err = service.ProfileRepository.AddProfile(ctx, &profile)
	if err != nil {
		return nil, err
	}

	if err := service.NotificationApproachService.EnableUserNotificationForUser(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to enable notification for user")
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

func (service *ProfileService) UpdateProfile(
	ctx context.Context,
	userID uint,
	profile *UserProfile,
) error {
	user, err := service.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return UserNotFound
	}

	checkedTypes := []NotificationType{}
	for _, approach := range profile.NotificationApproaches {
		if !approach.Enabled {
			continue
		}
		checkedTypes = append(checkedTypes, approach.Name)
	}

	if err := service.validateProfile(profile, user, checkedTypes); err != nil {
		return err
	}

	err = service.Engine.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := service.ProfileRepository.TransactionalUpdateProfile(tx, userID, profile); err != nil {
			return err
		}
		// TODO: update notification approaches in db
		if err := service.NotificationApproachService.TransactionUpdateNotificationApproaches(tx, profile.NotificationApproaches); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil

}

func (service *ProfileService) validateProfile(profile *UserProfile, user *User, checkedNotificationTypes []NotificationType) error {
	notificationTypeSet := SetFromSlice(checkedNotificationTypes)
	if notificationTypeSet.Contains(NotificationTypeEmail) && !user.IsVerified {
		return internal.UserNotVerified
	}
	if notificationTypeSet.Contains(NotificationTypeSMS) {
		if len(profile.PhoneNumber) == 0 {
			return internal.ProfilePhoneNumberRequired
		}
		if !profile.IsPhoneNumberVerified {
			return internal.ProfilePhoneNumberNotVerified
		}
	}
	if notificationTypeSet.Contains(NotificationTypeLineMessage) && user.Platform.Name != string(PlatformLine) {
		return internal.UserPlatformNotLinePlatform
	}

	return nil
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
	if err != nil {
		return nil, nil, err
	}
	if profile == nil {
		return nil, nil, internal.ProfileNotCreated
	}

	if profile.HasProfilePicture() {
		log.Info().Msgf("Deleting old profile image %s", profile.ProfilePictureObjectName)
		if err := service.FileUploadService.DeleteFile(ctx, profile.ProfilePictureObjectName); err != nil {
			log.Warn().Err(err).Msg("Failed to delete old profile image")
		}
	}

	objectName := uuid.New().String()

	uploadInfo, err := service.FileUploadService.UploadFile(ctx, objectName, file)
	if err != nil {
		return nil, nil, err
	}
	profile.ProfilePictureObjectName = objectName

	if err := service.UpdateProfile(ctx, userID, &UserProfile{ProfilePictureObjectName: uploadInfo.Key}); err != nil {
		return nil, nil, err
	}

	stringUrl, err := service.FileUploadService.GetFileExpiresIn(ctx, uploadInfo.Key, 5*time.Minute)
	if err != nil {
		return nil, nil, err
	}
	url := URL(stringUrl)
	return &url, uploadInfo, nil
}

func (service *ProfileService) GetProfileImage(ctx context.Context, userID uint) (*URL, error) {
	profile, err := service.FindProfileByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	url, err := service.FileUploadService.GetFileExpiresIn(ctx, profile.ProfilePictureObjectName, 30*time.Minute)
	if err != nil {
		return nil, err
	}
	return (*URL)(&url), nil

}

func (service *ProfileService) GetAllCategories(ctx context.Context) ([]*ProductCategory, error) {
	return service.ProfileRepository.FindAllCategories(ctx)
}
