package repository

import (
	"context"
	"gorm.io/gorm"
)

type IProfileRepository interface {
	FindProfileByUserId(ctx context.Context, customerId uint) (*UserProfile, error)
	AddProfile(ctx context.Context, profile *UserProfile) error
	UpdateProfile(ctx context.Context, profile *UserProfile, values any) error
	FindProfileByID(ctx context.Context, ID uint) (*UserProfile, error)
}

type ProfileRepository struct {
	Engine *gorm.DB
}

func (repo ProfileRepository) FindProfileByID(ctx context.Context, ID uint) (*UserProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (repo ProfileRepository) FindProfileByUserId(ctx context.Context, customerId uint) (*UserProfile, error) {
	profile := UserProfile{}
	tx := repo.Engine.WithContext(ctx).Where("customer_id = ?", customerId).First(&profile)
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
