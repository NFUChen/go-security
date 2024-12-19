package service

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	. "go-security/erp/internal/repository"
	"gorm.io/gorm"
)

type NotificationApproachService struct {
	ProfileService                 *ProfileService
	NotificationApproachRepository INotificationApproachRepository
}

func NewNotificationApproachService(notificationApproachRepository INotificationApproachRepository) *NotificationApproachService {
	return &NotificationApproachService{
		NotificationApproachRepository: notificationApproachRepository,
	}
}

func (service *NotificationApproachService) InjectProfileService(profileService *ProfileService) {
	service.ProfileService = profileService
}

func (service *NotificationApproachService) GetAllNotificationTypes() []NotificationType {
	return []NotificationType{
		NotificationTypeEmail,
		NotificationTypeSMS,
		NotificationTypeLineMessage,
	}
}

func (service *NotificationApproachService) createDefaultNotificationApproaches(userID uint) []*NotificationApproach {
	notificationTypes := service.GetAllNotificationTypes()

	var approaches []*NotificationApproach
	for _, notificationType := range notificationTypes {
		approach := &NotificationApproach{
			UserID:  userID,
			Name:    notificationType,
			Enabled: false,
		}
		approaches = append(approaches, approach)
	}
	return approaches
}

func (service *NotificationApproachService) IsUserNotificationEnabled(ctx context.Context, userID uint) bool {
	numberOfApproaches, err := service.NotificationApproachRepository.GetNumberOfApproachesByUserID(ctx, userID)
	if err != nil {
		return false
	}
	numberOfAvailableApproaches := len(service.GetAllNotificationTypes())
	isCorrupted := (numberOfApproaches != 0) && (numberOfApproaches < numberOfAvailableApproaches)
	if isCorrupted {
		log.Warn().Msg("Notification approaches are corrupted, resetting...")
		if err := service.handleResetApproaches(ctx, userID); err != nil {
			log.Error().Err(err).Msg("Failed to reset notification approaches")
			return false
		}
		return false
	}
	return numberOfApproaches == numberOfAvailableApproaches
}

func (service *NotificationApproachService) TransactionUpdateNotificationApproaches(tx *gorm.DB, approaches []NotificationApproach) error {
	ptrApproaches := []*NotificationApproach{}
	for _, approach := range approaches {
		ptrApproaches = append(ptrApproaches, &approach)
	}
	return service.NotificationApproachRepository.TransactionalUpdateNotificationApproaches(tx, ptrApproaches)
}

func (service *NotificationApproachService) UpdateNotificationApproaches(ctx context.Context, approaches []*NotificationApproach) error {
	return service.NotificationApproachRepository.UpdateNotificationApproaches(ctx, approaches)
}

func (service *NotificationApproachService) handleResetApproaches(ctx context.Context, userID uint) error {
	defaultApproaches := service.createDefaultNotificationApproaches(userID)
	if err := service.NotificationApproachRepository.ResetNotificationApproaches(ctx, userID, defaultApproaches); err != nil {
		log.Warn().Err(err).Msg("Failed to reset notification approaches...")
		return err
	}
	return nil
}

func (service *NotificationApproachService) EnableUserNotificationForUser(ctx context.Context, userID uint) error {
	if !service.ProfileService.IsProfileExists(ctx, userID) {
		log.Warn().Msgf("User %d: Profile not exists, cannot enable notification, create a default one", userID)
		if _, err := service.ProfileService.CreateDefaultProfile(ctx, userID); err != nil {
			return err
		}
	}
	number, err := service.NotificationApproachRepository.GetNumberOfApproachesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get number of approaches: %w", err)
	}
	if number == len(service.GetAllNotificationTypes()) {
		log.Info().Msgf("User %d: Notification already enabled...", userID)
		return nil
	}
	approaches := service.createDefaultNotificationApproaches(userID)

	if err := service.NotificationApproachRepository.SaveNotificationApproaches(ctx, approaches); err != nil {
		return err
	}
	return nil

}
