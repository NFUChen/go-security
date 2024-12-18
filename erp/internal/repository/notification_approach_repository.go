package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type INotificationApproachRepository interface {
	UpdateNotificationApproaches(ctx context.Context, approaches []*NotificationApproach) error
	TransactionalUpdateNotificationApproaches(tx *gorm.DB, approaches []*NotificationApproach) error
	GetNumberOfApproachesByUserID(ctx context.Context, userID uint) (int, error)
	SaveNotificationApproaches(ctx context.Context, approaches []*NotificationApproach) error
	ResetNotificationApproaches(ctx context.Context, userID uint, resetWithApproaches []*NotificationApproach) error
}

type NotificationApproachRepository struct {
	Engine *gorm.DB
}

func (repo NotificationApproachRepository) TransactionalUpdateNotificationApproaches(tx *gorm.DB, approaches []*NotificationApproach) error {
	for _, approach := range approaches {
		log.Info().Msgf("Updating notification approach: %v", approach)
		// Perform upsert: update existing or insert new if it doesn't exist
		existingApproach := &NotificationApproach{}
		if err := tx.First(existingApproach, "user_id = ? AND name = ?", approach.UserID, approach.Name).Error; err != nil {
			log.Error().Err(err).Msg("Failed to find existing notification approach")
		}
		if err := tx.Model(existingApproach).Update("enabled", approach.Enabled).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update notification approach status")
			return err
		}
	}
	return nil
}

func (repo NotificationApproachRepository) ResetNotificationApproaches(ctx context.Context, userID uint, resetWithApproaches []*NotificationApproach) error {
	return repo.Engine.WithContext(ctx).Delete(&NotificationApproach{}, "user_id = ?", userID).Error
}

func (repo NotificationApproachRepository) GetNumberOfApproachesByUserID(ctx context.Context, userID uint) (int, error) {
	var count int64
	if err := repo.Engine.WithContext(ctx).Model(&NotificationApproach{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count notification approaches")
		return 0, err
	}
	return int(count), nil
}

func (repo NotificationApproachRepository) UpdateNotificationApproaches(ctx context.Context, approaches []*NotificationApproach) error {
	return repo.Engine.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			for _, approach := range approaches {
				log.Info().Msgf("Updating notification approach: %v", approach)
				if err := tx.Delete(&NotificationApproach{}, "user_id = ? AND name = ?", approach.UserID, approach.Name).Error; err != nil {
					log.Error().Err(err).Msg("Failed to delete existing notification approach")
					return err
				}
				if err := tx.Create(approach).Error; err != nil {
					log.Error().Err(err).Msg("Failed to create new notification approach")
					return err
				}
			}
			return nil
		})

}

func (repo NotificationApproachRepository) SaveNotificationApproaches(ctx context.Context, approaches []*NotificationApproach) error {
	return repo.Engine.WithContext(ctx).Create(approaches).Error
}

func NewNotificationApproachRepository(engine *gorm.DB) *NotificationApproachRepository {
	return &NotificationApproachRepository{
		Engine: engine,
	}
}
