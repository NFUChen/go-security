package service

import (
	"context"
	"github.com/rs/zerolog/log"
	. "go-security/erp/internal/repository"
)

type NotificationApproachService struct {
	NotificationApproachRepository INotificationApproachRepository
}

func NewNotificationApproachService(notificationApproachRepository INotificationApproachRepository) *NotificationApproachService {
	return &NotificationApproachService{
		NotificationApproachRepository: notificationApproachRepository,
	}
}

func (service *NotificationApproachService) GetAllNotificationTypes() []NotificationType {
	return []NotificationType{
		NotificationTypeEmail,
		NotificationTypeSMS,
		NotificationTypeLineMessage,
	}
}

func (service *NotificationApproachService) createDefaultNotificationApproaches() []NotificationApproach {
	notificationTypes := service.GetAllNotificationTypes()

	var approaches []NotificationApproach
	for _, notificationType := range notificationTypes {
		approach := NotificationApproach{
			Name:    notificationType,
			Enabled: false,
		}
		approaches = append(approaches, approach)
	}
	return approaches
}

func (service *NotificationApproachService) IsUserNotificationEnabled(ctx context.Context, userID uint) bool {
	numberOfApproaches := service.NotificationApproachRepository.GetNumberOfApproachesByUserID(ctx, userID)
	numberOfAvailableApproaches := len(service.GetAllNotificationTypes())
	return numberOfApproaches == numberOfAvailableApproaches
}

func (service *NotificationApproachService) UpdateNotificationApproaches(ctx context.Context, userID uint, approaches []NotificationApproach) error {
	return service.NotificationApproachRepository.UpdateNotificationApproaches(ctx, userID, approaches)
}

func (service *NotificationApproachService) EnableUserNotificationForUser(ctx context.Context, userID uint) error {
	number := service.NotificationApproachRepository.GetNumberOfApproachesByUserID(ctx, userID)
	if number == len(service.GetAllNotificationTypes()) {
		log.Info().Msgf("User %d: Notification already enabled...", userID)
		return nil
	}
	approaches := service.createDefaultNotificationApproaches()

	if err := service.NotificationApproachRepository.SaveNotificationApproaches(ctx, userID, approaches); err != nil {
		return err
	}
	return nil

}
