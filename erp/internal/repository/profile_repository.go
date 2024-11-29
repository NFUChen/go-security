package repository

import (
	"context"
	"gorm.io/gorm"
)

type IProfileRepository interface {
	FindProfileByUserID(ctx context.Context, customerId uint) (*UserProfile, error)
	FindProfileByID(ctx context.Context, userID uint) (*UserProfile, error)
	AddProfile(ctx context.Context, profile *UserProfile) error
	UpdateProfile(ctx context.Context, profile *UserProfile, values any) error
	FindAllProfiles(ctx context.Context) ([]*UserProfile, error)

	TransactionalUpdateProfile(tx *gorm.DB, profile *UserProfile, values any) error
	TransactionalAddProfile(tx *gorm.DB, profile *UserProfile) error
}

type ProfileRepository struct {
	Engine *gorm.DB
}

func (repo *ProfileRepository) TransactionalUpdateProfile(session *gorm.DB, profile *UserProfile, values any) error {
	return session.Model(&profile).Updates(values).Error
}

func (repo *ProfileRepository) TransactionalAddProfile(tx *gorm.DB, profile *UserProfile) error {
	return tx.Create(profile).Error
}

func (repo *ProfileRepository) createPreloadQuery(ctx context.Context) *gorm.DB {

	preloadsRequired := []string{
		"NotificationApproaches", "PricingPolicy",
	}
	engine := repo.Engine.WithContext(ctx)
	for _, preload := range preloadsRequired {
		engine = engine.Preload(preload)
	}

	return engine
}

func (repo *ProfileRepository) FindAllProfiles(ctx context.Context) ([]*UserProfile, error) {
	profiles := []*UserProfile{}
	tx := repo.createPreloadQuery(ctx).Find(&profiles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return profiles, nil
}

func (repo ProfileRepository) FindProfileByID(ctx context.Context, ID uint) (*UserProfile, error) {
	profile := UserProfile{}
	tx := repo.createPreloadQuery(ctx).Where("id = ?", ID).First(&profile)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &profile, nil
}

func (repo ProfileRepository) FindProfileByUserID(ctx context.Context, ID uint) (*UserProfile, error) {
	profile := UserProfile{}
	tx := repo.createPreloadQuery(ctx).First(&profile, ID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &profile, nil
}

func (repo ProfileRepository) AddProfile(ctx context.Context, profile *UserProfile) error {
	return repo.Engine.WithContext(ctx).Create(profile).Error
}

func (repo ProfileRepository) UpdateProfile(ctx context.Context, profile *UserProfile, values any) error {
	tx := repo.Engine.WithContext(ctx).Model(&profile).Updates(values)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewProfileRepository(engine *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		Engine: engine,
	}
}
