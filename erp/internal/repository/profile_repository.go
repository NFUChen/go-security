package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type IProfileRepository interface {
	FindProfileByUserID(ctx context.Context, userID uint) (*UserProfile, error)
	FindProfileByID(ctx context.Context, userID uint) (*UserProfile, error)
	AddProfile(ctx context.Context, profile *UserProfile) error
	UpdateProfile(ctx context.Context, userID uint, withProfile *UserProfile) error
	FindAllProfiles(ctx context.Context) ([]*UserProfile, error)

	TransactionalUpdateProfile(tx *gorm.DB, userID uint, profile *UserProfile) error
	TransactionalAddProfile(tx *gorm.DB, profile *UserProfile) error
	FindAllCategories(ctx context.Context) ([]*ProductCategory, error)
}

type ProfileRepository struct {
	Engine                         *gorm.DB
	NotificationApproachRepository INotificationApproachRepository
}

func NewProfileRepository(engine *gorm.DB, notificationApproachRepository INotificationApproachRepository) *ProfileRepository {
	return &ProfileRepository{
		Engine:                         engine,
		NotificationApproachRepository: notificationApproachRepository,
	}
}

func (repo *ProfileRepository) TransactionalUpdateProfile(session *gorm.DB, userID uint, updatedProfile *UserProfile) error {
	// Find the existing profile first
	var existingProfile UserProfile
	tx := session.Where("user_id = ?", userID).First(&existingProfile)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("profile with user ID %d not found", userID)
		}
		return tx.Error
	}

	// Update only the fields you need
	tx = session.Model(&existingProfile).Updates(updatedProfile)
	return tx.Error
}

func (repo *ProfileRepository) TransactionalAddProfile(tx *gorm.DB, profile *UserProfile) error {
	return tx.Create(profile).Error
}

func (repo *ProfileRepository) createProductPreloadQuery(ctx context.Context) *gorm.DB {
	engine := repo.Engine.WithContext(ctx).Preload("PricingPolicy")
	return engine
}

func (repo *ProfileRepository) FindAllProfiles(ctx context.Context) ([]*UserProfile, error) {
	profiles := []*UserProfile{}
	tx := repo.createProductPreloadQuery(ctx).Find(&profiles)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, profile := range profiles {
		approaches, err := repo.NotificationApproachRepository.GetNotificationApproachesByUserID(ctx, profile.UserID)
		if err != nil {
			return nil, err
		}
		profile.NotificationApproaches = approaches
	}

	return profiles, nil
}

func (repo *ProfileRepository) FindProfileByID(ctx context.Context, ID uint) (*UserProfile, error) {
	profile := UserProfile{}
	tx := repo.createProductPreloadQuery(ctx).Where("id = ?", ID).First(&profile)
	if tx.Error != nil {
		return nil, tx.Error
	}

	approaches, err := repo.NotificationApproachRepository.GetNotificationApproachesByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}
	profile.NotificationApproaches = approaches

	return &profile, nil
}

func (repo *ProfileRepository) FindProfileByUserID(ctx context.Context, userID uint) (*UserProfile, error) {
	profile := UserProfile{}
	tx := repo.createProductPreloadQuery(ctx).Where("user_id = ?", userID).First(&profile)
	if tx.Error != nil {
		return nil, tx.Error
	}

	approaches, err := repo.NotificationApproachRepository.GetNotificationApproachesByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}
	profile.NotificationApproaches = approaches

	return &profile, nil
}

func (repo *ProfileRepository) AddProfile(ctx context.Context, profile *UserProfile) error {
	return repo.Engine.WithContext(ctx).Create(profile).Error
}

func (repo *ProfileRepository) UpdateProfile(ctx context.Context, userID uint, withProfile *UserProfile) error {
	var profile UserProfile

	// Find the profile with preloading
	tx := repo.createProductPreloadQuery(ctx).Where("user_id = ?", userID).First(&profile)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("profile with user ID %d not found", userID)
		}
		return tx.Error
	}

	// Perform withProfile in a transaction
	err := repo.Engine.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&profile).Updates(withProfile).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (repo *ProfileRepository) FindAllCategories(ctx context.Context) ([]*ProductCategory, error) {
	categories := []*ProductCategory{}
	tx := repo.Engine.WithContext(ctx).Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}
